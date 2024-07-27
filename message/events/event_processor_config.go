package events

import (
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

var Marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
}

func NewEventProcessorConfig(rdb *redis.Client, logger watermill.LoggerAdapter) cqrs.EventProcessorConfig {
	return cqrs.EventProcessorConfig{
		SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(
				redisstream.SubscriberConfig{
					Client:        rdb,
					ConsumerGroup: "svc-tickets.events." + params.HandlerName,
				},
				logger,
			)
		},
		GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			handlerEvent := params.EventHandler.NewEvent()
			event, ok := handlerEvent.(entities.Event)
			if !ok {
				return "", fmt.Errorf("invalid event type: %T doesn't implement entities.Event", handlerEvent)
			}

			if event.IsInternal() {
				return "internal-events.svc-tickets." + params.EventName, nil
			}

			return "events." + params.EventName, nil
		},
		Marshaler: Marshaler,
		Logger:    logger,
	}
}
