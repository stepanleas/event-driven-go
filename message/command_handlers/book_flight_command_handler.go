package command_handlers

import (
	"context"
	"fmt"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type BookFlightCommandHandler struct {
	transportationClient contracts.TransportationService
	eventBus             *cqrs.EventBus
}

func NewBookFlightCommandHandler(transportationService contracts.TransportationService, eventBus *cqrs.EventBus) BookFlightCommandHandler {
	return BookFlightCommandHandler{
		transportationClient: transportationService,
		eventBus:             eventBus,
	}
}

func (h BookFlightCommandHandler) Handle(ctx context.Context, command *entities.BookFlight) error {
	resp, err := h.transportationClient.BookFlight(ctx, entities.BookFlightTicketRequest{
		CustomerEmail:  command.CustomerEmail,
		FlightID:       command.FlightID,
		PassengerNames: command.Passengers,
		ReferenceId:    command.ReferenceID,
		IdempotencyKey: command.IdempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to book flight: %w", err)
	}

	err = h.eventBus.Publish(ctx, entities.FlightBooked_v1{
		Header:      entities.NewEventHeader(),
		FlightID:    command.FlightID,
		TicketIDs:   resp.TicketIds,
		ReferenceID: command.ReferenceID,
	})
	if err != nil {
		return fmt.Errorf("failed to publish FlightBooked_v1 event: %w", err)
	}

	return nil
}
