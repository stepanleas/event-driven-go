package message

import (
	"fmt"
	"tickets/entities"
	"tickets/message/contracts"
	"tickets/message/events"
	"tickets/message/events/outbox"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewWatermillRouter(
	receiptsService contracts.ReceiptsService,
	spreadsheetsService contracts.SpreadsheetsAPI,
	dataLake contracts.DataLake,
	postgresSubscriber message.Subscriber,
	redisPublisher message.Publisher,
	redisSubscriber message.Subscriber,
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

	router.AddNoPublisherHandler(
		"events_splitter",
		"events",
		redisSubscriber,
		func(msg *message.Message) error {
			eventName := events.Marshaler.NameFromMessage(msg)
			if eventName == "" {
				return fmt.Errorf("cannot get event name from message")
			}

			return redisPublisher.Publish("events."+eventName, msg)
		},
	)

	router.AddNoPublisherHandler(
		"events_store",
		"events",
		redisSubscriber,
		func(msg *message.Message) error {
			eventName := events.Marshaler.NameFromMessage(msg)
			if eventName == "" {
				return fmt.Errorf("cannot get event name from message")
			}

			type Event struct {
				Header entities.EventHeader `json:"header"`
			}

			var event Event
			if err := events.Marshaler.Unmarshal(msg, &event); err != nil {
				return fmt.Errorf("cannot unmarshal event: %w", err)
			}

			return dataLake.Store(msg.Context(), entities.DataLakeEvent{
				EventID:      event.Header.ID,
				PublishedAt:  time.Now(),
				EventName:    eventName,
				EventPayload: msg.Payload,
			})
		},
	)

	outbox.AddForwarderHandler(postgresSubscriber, redisPublisher, router, logger)

	return router
}
