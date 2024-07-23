package events

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewEventBus(pub message.Publisher) *cqrs.EventBus {
	eventBus, err := cqrs.NewEventBusWithConfig(
		pub,
		cqrs.EventBusConfig{
			GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
				return fmt.Sprintf("events.%s", params.EventName), nil
			},
			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	return eventBus
}
