package service

import (
	"context"
	"fmt"
	stdHTTP "net/http"
	"tickets/db"
	"tickets/db/read_model"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/command_handlers/contract"
	"tickets/message/commands"
	"tickets/message/contracts"
	"tickets/message/events"
	"tickets/message/events/outbox"
	"tickets/migrations"
	"tickets/observability"
	"tickets/process_manager"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/sdk/trace"

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
	dataLake        contracts.DataLake
	opsReadModel    read_model.OpsBookingReadModel
	tracerProvider  *trace.TracerProvider
}

type ReceiptService interface {
	contracts.ReceiptsService
	contract.ReceiptsService
}

func New(
	dbConn *sqlx.DB,
	redisClient *redis.Client,
	spreadsheetsService contracts.SpreadsheetsAPI,
	receiptsService ReceiptService,
	filesAPI contracts.FilesAPI,
	deadNationAPI contracts.DeadNationApi,
	paymentsService contract.PaymentsService,
) Service {
	tracerProvider := observability.ConfigureTracerProvider()

	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	redisPublisher := message.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}
	redisPublisher = observability.TracingPublisherDecorator{Publisher: redisPublisher}

	redisSubscriber := message.NewRedisSubscriber(redisClient, watermillLogger)
	eventBus := events.NewEventBus(redisPublisher)

	ticketsRepo := db.NewTicketRepository(dbConn)
	showRepo := db.NewShowRepository(dbConn)
	bookingRepo := db.NewBookingRepository(dbConn)
	dataLake := db.NewDataLake(dbConn)
	opsReadModel := read_model.NewOpsBookingReadModel(dbConn, eventBus)

	postgresSubscriber := outbox.NewPostgresSubscriber(dbConn.DB, watermillLogger)

	watermillRouter := message.NewWatermillRouter(
		dataLake,
		postgresSubscriber,
		redisPublisher,
		redisSubscriber,
		watermillLogger,
	)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		watermillRouter,
		events.NewEventProcessorConfig(redisClient, watermillLogger),
	)
	if err != nil {
		panic(err)
	}

	commandBus := commands.NewCommandBus(redisPublisher)
	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		watermillRouter,
		commands.NewCommandProcessorConfig(redisClient, watermillLogger),
	)
	if err != nil {
		panic(err)
	}

	commands.AddCommandProcessorHandlers(commandProcessor, eventBus, bookingRepo, receiptsService, paymentsService)

	vipBundleRepo := db.NewVipBundleRepository(dbConn)
	vipBundlePM := process_manager.NewVipBundleProcessManager(commandBus, eventBus, vipBundleRepo)

	events.AddEventProcessorHandlers(
		eventProcessor,
		eventBus,
		receiptsService,
		spreadsheetsService,
		ticketsRepo,
		showRepo,
		filesAPI,
		deadNationAPI,
		opsReadModel,
		vipBundlePM,
	)

	echoRouter := ticketsHttp.NewHttpRouter(
		eventBus,
		commandBus,
		spreadsheetsService,
		ticketsRepo,
		showRepo,
		bookingRepo,
		vipBundleRepo,
		opsReadModel,
	)

	return Service{
		db:              dbConn,
		watermillRouter: watermillRouter,
		echoRouter:      echoRouter,
		dataLake:        dataLake,
		opsReadModel:    opsReadModel,
		tracerProvider:  tracerProvider,
	}
}

func (s Service) Run(ctx context.Context) error {
	if err := db.InitializeDatabaseSchema(s.db); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	go func() {
		if err := migrations.MigrateReadModel(ctx, s.dataLake, s.opsReadModel); err != nil {
			log.FromContext(ctx).Errorf("failed to migrate read model: %v", err)
		}
	}()

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

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.tracerProvider.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
