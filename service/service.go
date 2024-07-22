package service

import (
	"context"
	"fmt"
	stdHTTP "net/http"
	"tickets/db"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/commands"
	"tickets/message/contracts"
	"tickets/message/events"
	"tickets/message/events/outbox"

	_ "github.com/lib/pq"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {
	log.Init(logrus.InfoLevel)
}

type Service struct {
	db              *sqlx.DB
	watermillRouter *watermillMessage.Router
	echoRouter      *echo.Echo
}

func New(
	dbConn *sqlx.DB,
	redisClient *redis.Client,
	spreadsheetsService contracts.SpreadsheetsAPI,
	receiptsService contracts.ReceiptsService,
	filesAPI contracts.FilesAPI,
	deadNationAPI contracts.DeadNationApi,
) Service {
	ticketsRepo := db.NewTicketRepository(dbConn)
	showRepo := db.NewShowRepository(dbConn)
	bookingRepo := db.NewBookingRepository(dbConn)

	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	var redisPublisher watermillMessage.Publisher
	redisPublisher = message.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}

	postgresSub := outbox.NewPostgresSubscriber(dbConn.DB, watermillLogger)

	watermillRouter := message.NewWatermillRouter(
		receiptsService,
		spreadsheetsService,
		postgresSub,
		redisPublisher,
		watermillLogger,
	)

	eventBus := events.NewEventBus(redisPublisher)
	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		watermillRouter,
		events.NewEventProcessorConfig(redisClient, watermillLogger),
	)
	if err != nil {
		panic(err)
	}

	events.AddEventProcessorHandlers(
		eventProcessor,
		eventBus,
		receiptsService,
		spreadsheetsService,
		ticketsRepo,
		showRepo,
		filesAPI,
		deadNationAPI,
	)

	commandBus := commands.NewCommandBus(redisPublisher)

	echoRouter := ticketsHttp.NewHttpRouter(
		eventBus,
		commandBus,
		spreadsheetsService,
		ticketsRepo,
		showRepo,
		bookingRepo,
	)

	return Service{
		db:              dbConn,
		watermillRouter: watermillRouter,
		echoRouter:      echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	if err := db.InitializeDatabaseSchema(s.db); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && err != stdHTTP.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
