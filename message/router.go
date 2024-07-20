package message

import (
	"tickets/message/contracts"
	"tickets/message/events/outbox"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewWatermillRouter(
	receiptsService contracts.ReceiptsService,
	spreadsheetsService contracts.SpreadsheetsAPI,
	postgresSubscriber message.Subscriber,
	publisher message.Publisher,
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

	outbox.AddForwarderHandler(postgresSubscriber, publisher, router, logger)

	return router
}
