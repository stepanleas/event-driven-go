package command_handlers

import (
	"context"
	"fmt"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type BookTaxiCommandHandler struct {
	transportationClient contracts.TransportationService
	eventBus             *cqrs.EventBus
}

func NewBookTaxiCommandHandler(transportationService contracts.TransportationService, eventBus *cqrs.EventBus) BookTaxiCommandHandler {
	return BookTaxiCommandHandler{
		transportationClient: transportationService,
		eventBus:             eventBus,
	}
}

func (h BookTaxiCommandHandler) Handle(ctx context.Context, command *entities.BookTaxi) error {
	resp, err := h.transportationClient.BookTaxi(ctx, entities.BookTaxiRequest{
		CustomerEmail:      command.CustomerEmail,
		NumberOfPassengers: command.NumberOfPassengers,
		PassengerName:      command.CustomerName,
		ReferenceId:        command.ReferenceID,
		IdempotencyKey:     command.IdempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to book taxi: %w", err)
	}

	err = h.eventBus.Publish(ctx, entities.TaxiBooked_v1{
		Header:        entities.NewEventHeader(),
		TaxiBookingID: resp.TaxiBookingId,
		ReferenceID:   command.ReferenceID,
	})
	if err != nil {
		return fmt.Errorf("failed to publish TaxiBooked_v1 event: %w", err)
	}

	return nil
}
