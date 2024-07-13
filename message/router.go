package message

import (
	"tickets/message/contracts"
	"tickets/message/handlers"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

const brokenMessageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

func NewWatermillRouter(
	receiptsService contracts.ReceiptsService,
	spreadsheetsService contracts.SpreadsheetsAPI,
	rdb *redis.Client,
	logger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(
		RetryMiddleware(logger).Middleware,
		CorrelationIDMiddleware,
		LoggingMiddleware,
	)

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

	router.AddNoPublisherHandler(
		"issue_receipt_handler",
		"TicketBookingConfirmed",
		issueReceiptSub,
		handlers.NewIssueReceiptsHandler(receiptsService).Handle,
	)

	router.AddNoPublisherHandler(
		"append_to_tracker_handler",
		"TicketBookingConfirmed",
		appendToTrackerSub,
		handlers.NewAppendToTrackerHandler(spreadsheetsService).Handle,
	)

	router.AddNoPublisherHandler(
		"tickets_to_refund_handler",
		"TicketBookingCanceled",
		cancelTicketSub,
		handlers.NewTicketsToRefundHandler(spreadsheetsService).Handle,
	)

	return router
}
