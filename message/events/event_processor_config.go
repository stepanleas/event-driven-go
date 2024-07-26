package events

import (
	"fmt"

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
			return fmt.Sprintf("events.%s", params.EventName), nil
		},
		Marshaler: Marshaler,
		Logger:    logger,
	}
}
