package commands

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

var marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
}

func NewCommandProcessorConfig(rdb *redis.Client, logger watermill.LoggerAdapter) cqrs.CommandProcessorConfig {
	return cqrs.CommandProcessorConfig{
		SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(
				redisstream.SubscriberConfig{
					Client:        rdb,
					ConsumerGroup: "svc-tickets.commands." + params.HandlerName,
				},
				logger,
			)
		},
		GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
			return params.CommandName, nil
		},
		Marshaler: marshaler,
		Logger:    logger,
	}
}
