package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/receipts"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/spreadsheets"
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const brokenMessageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type TicketsStatusRequest struct {
	Tickets []TicketStatus `json:"tickets"`
}

type TicketStatus struct {
	TicketID      string `json:"ticket_id"`
	Status        string `json:"status"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type EventHeader struct {
	ID          string `json:"id"`
	PublishedAt string `json:"published_at"`
}

func NewEventHeader() EventHeader {
	return EventHeader{
		ID:          watermill.NewUUID(),
		PublishedAt: time.Now().Format(time.RFC3339),
	}
}

type TicketBookingConfirmed struct {
	Header        EventHeader `json:"header"`
	TicketID      string      `json:"ticket_id"`
	CustomerEmail string      `json:"customer_email"`
	Price         Money       `json:"price"`
}

type TicketBookingCanceled struct {
	Header        EventHeader `json:"header"`
	TicketID      string      `json:"ticket_id"`
	CustomerEmail string      `json:"customer_email"`
	Price         Money       `json:"price"`
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

	receiptsClient := NewReceiptsClient(clients)
	spreadsheetsClient := NewSpreadsheetsClient(clients)

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
				event := TicketBookingConfirmed{
					Header:        NewEventHeader(),
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
				event := TicketBookingCanceled{
					Header:        NewEventHeader(),
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

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(RetryMiddleware(logger).Middleware)
	router.AddMiddleware(CorrelationIDMiddleware)
	router.AddMiddleware(LoggingMiddleware)

	router.AddNoPublisherHandler(
		"issue_receipt_handler",
		"TicketBookingConfirmed",
		issueReceiptSub,
		func(msg *message.Message) error {
			if msg.UUID == brokenMessageID {
				return nil
			}

			if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
				return nil
			}

			var payload TicketBookingConfirmed
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				return err
			}

			return receiptsClient.IssueReceipt(msg.Context(), IssueReceiptRequest{
				TicketID: payload.TicketID,
				Price:    payload.Price,
			})
		},
	)

	router.AddNoPublisherHandler(
		"append_to_tracker_handler",
		"TicketBookingConfirmed",
		appendToTrackerSub,
		func(msg *message.Message) error {
			if msg.UUID == brokenMessageID {
				return nil
			}

			if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
				return nil
			}

			var payload TicketBookingConfirmed
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				return err
			}

			return spreadsheetsClient.AppendRow(
				msg.Context(),
				"tickets-to-print",
				[]string{payload.TicketID, payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency},
			)
		},
	)

	router.AddNoPublisherHandler(
		"tickets_to_refund_handler",
		"TicketBookingCanceled",
		cancelTicketSub,
		func(msg *message.Message) error {
			if msg.UUID == brokenMessageID {
				return nil
			}

			if msg.Metadata.Get("type") != "TicketBookingCanceled" {
				return nil
			}

			var payload TicketBookingCanceled
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				return err
			}

			return spreadsheetsClient.AppendRow(
				msg.Context(),
				"tickets-to-refund",
				[]string{payload.TicketID, payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency},
			)
		},
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

func CorrelationIDMiddleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		ctx := msg.Context()

		reqCorrelationId := msg.Metadata.Get("correlation_id")
		if reqCorrelationId == "" {
			reqCorrelationId = watermill.NewUUID()
		}

		ctx = log.ToContext(ctx, logrus.WithFields(logrus.Fields{"correlation_id": reqCorrelationId}))
		ctx = log.ContextWithCorrelationID(ctx, reqCorrelationId)

		msg.SetContext(ctx)

		return h(msg)
	}
}

func LoggingMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := log.FromContext(msg.Context()).WithField("message_uuid", msg.UUID)

		logger.Info("Handling a message")

		msgs, err := next(msg)
		if err != nil {
			logger.WithError(err).Error("Message handling error")
		}

		return msgs, err
	}
}

func RetryMiddleware(logger watermill.LoggerAdapter) middleware.Retry {
	return middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          logger,
	}
}

type ReceiptsClient struct {
	clients *clients.Clients
}

func NewReceiptsClient(clients *clients.Clients) ReceiptsClient {
	return ReceiptsClient{
		clients: clients,
	}
}

func (c ReceiptsClient) IssueReceipt(ctx context.Context, request IssueReceiptRequest) error {
	body := receipts.PutReceiptsJSONRequestBody{
		TicketId: request.TicketID,
		Price: receipts.Money{
			MoneyAmount:   request.Price.Amount,
			MoneyCurrency: request.Price.Currency,
		},
	}

	receiptsResp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if receiptsResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", receiptsResp.StatusCode())
	}

	return nil
}

type SpreadsheetsClient struct {
	clients *clients.Clients
}

func NewSpreadsheetsClient(clients *clients.Clients) SpreadsheetsClient {
	return SpreadsheetsClient{
		clients: clients,
	}
}

func (c SpreadsheetsClient) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	request := spreadsheets.PostSheetsSheetRowsJSONRequestBody{
		Columns: row,
	}

	sheetsResp, err := c.clients.Spreadsheets.PostSheetsSheetRowsWithResponse(ctx, spreadsheetName, request)
	if err != nil {
		return err
	}
	if sheetsResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", sheetsResp.StatusCode())
	}

	return nil
}
