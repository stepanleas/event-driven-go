package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"tickets/api"
	"tickets/events"
	"tickets/events/handlers"
	"tickets/valueobject"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const brokenMessageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

type TicketsStatusRequest struct {
	Tickets []TicketStatus `json:"tickets"`
}

type TicketStatus struct {
	TicketID      string            `json:"ticket_id"`
	Status        string            `json:"status"`
	CustomerEmail string            `json:"customer_email"`
	Price         valueobject.Money `json:"price"`
}

func main() {
	log.Init(logrus.InfoLevel)

	clients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	receiptsClient := api.NewReceiptsClient(clients)
	spreadsheetsClient := api.NewSpreadsheetsClient(clients)

	e := commonHTTP.NewEcho()

	logger := watermill.NewStdLogger(false, false)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, logger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, logger)
	if err != nil {
		panic(err)
	}

	cancelTicketSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "cancel-ticket",
	}, logger)
	if err != nil {
		panic(err)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/tickets-status", func(c echo.Context) error {
		var request TicketsStatusRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		for _, ticket := range request.Tickets {
			if ticket.Status == "confirmed" {
				event := events.TicketBookingConfirmed{
					Header:        events.NewEventHeader(),
					TicketID:      ticket.TicketID,
					CustomerEmail: ticket.CustomerEmail,
					Price:         ticket.Price,
				}

				payload, err := json.Marshal(event)
				if err != nil {
					return err
				}

				msg := message.NewMessage(watermill.NewUUID(), []byte(payload))
				msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-Id"))
				msg.Metadata.Set("type", "TicketBookingConfirmed")

				if err := pub.Publish("TicketBookingConfirmed", msg); err != nil {
					return err
				}
			} else if ticket.Status == "canceled" {
				event := events.TicketBookingCanceled{
					Header:        events.NewEventHeader(),
					TicketID:      ticket.TicketID,
					CustomerEmail: ticket.CustomerEmail,
					Price:         ticket.Price,
				}

				payload, err := json.Marshal(event)
				if err != nil {
					return err
				}

				msg := message.NewMessage(watermill.NewUUID(), []byte(payload))
				msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-Id"))
				msg.Metadata.Set("type", "TicketBookingCanceled")

				if err := pub.Publish("TicketBookingCanceled", msg); err != nil {
					return err
				}
			}

		}

		return c.NoContent(http.StatusOK)
	})

	router := events.NewRouter(logger)

	router.AddMiddleware(events.RetryMiddleware(logger).Middleware)
	router.AddMiddleware(events.CorrelationIDMiddleware)
	router.AddMiddleware(events.LoggingMiddleware)

	router.AddNoPublisherHandler(
		"issue_receipt_handler",
		"TicketBookingConfirmed",
		issueReceiptSub,
		handlers.NewIssueReceiptsHandler(receiptsClient).Handle,
	)

	router.AddNoPublisherHandler(
		"append_to_tracker_handler",
		"TicketBookingConfirmed",
		appendToTrackerSub,
		handlers.NewAppendToTrackerHandler(spreadsheetsClient).Handle,
	)

	router.AddNoPublisherHandler(
		"tickets_to_refund_handler",
		"TicketBookingCanceled",
		cancelTicketSub,
		handlers.NewTicketsToRefundHandler(spreadsheetsClient).Handle,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return router.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-router.Running()

		err := e.Start(":8080")

		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return e.Shutdown(context.Background())
	})

	if err := errgrp.Wait(); err != nil {
		panic(err)
	}
}
