package command_handlers

import (
	"context"
	"errors"
	"fmt"
	"tickets/db"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type BookShowTicketsCommandHandler struct {
	bookingRepo contracts.BookingRepository
	eventBus    *cqrs.EventBus
}

func NewBookShowTicketsCommandHandler(bookingRepo contracts.BookingRepository, eventBus *cqrs.EventBus) BookShowTicketsCommandHandler {
	return BookShowTicketsCommandHandler{
		bookingRepo: bookingRepo,
		eventBus:    eventBus,
	}
}

func (h BookShowTicketsCommandHandler) Handle(ctx context.Context, command *entities.BookShowTickets) error {
	err := h.bookingRepo.Add(ctx, entities.Booking{
		BookingID:       command.BookingID,
		ShowID:          command.ShowId,
		NumberOfTickets: command.NumberOfTickets,
		CustomerEmail:   command.CustomerEmail,
	})
	if errors.Is(err, db.ErrBookingAlreadyExists) {
		// now AddBooking is called via Pub/Sub, we are taking into account at-least-once delivery
		return nil
	}

	if errors.Is(err, db.ErrNoPlacesLeft) {
		publishErr := h.eventBus.Publish(ctx, entities.BookingFailed_v1{
			Header:        entities.NewEventHeader(),
			BookingID:     command.BookingID,
			FailureReason: err.Error(),
		})
		if publishErr != nil {
			return fmt.Errorf("failed to publish BookingFailed_v1 event: %w", publishErr)
		}
	}

	return err
}
