package message

import (
	"fmt"
	"tickets/entities"
	"tickets/message/contracts"
	"tickets/message/events"
	"tickets/message/events/outbox"
	"tickets/message/middleware"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/prometheus/client_golang/prometheus"
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
		middleware.RetryMiddleware(logger).Middleware,
		middleware.CorrelationIDMiddleware,
		middleware.LoggingMiddleware,
	)

	router.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) (events []*message.Message, err error) {
			start := time.Now()

			topic := message.SubscribeTopicFromCtx(msg.Context())
			handler := message.HandlerNameFromCtx(msg.Context())

			labels := prometheus.Labels{"topic": topic, "handler": handler}

			middleware.MessagesProcessedCounter.With(labels).Inc()

			msgs, err := h(msg)
			if err != nil {
				middleware.MessagesProcessingFailedCounter.With(labels).Inc()
			}

			middleware.MessagesProcessingDuration.With(labels).Observe(time.Since(start).Seconds())

			return msgs, err
		}
	})

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
