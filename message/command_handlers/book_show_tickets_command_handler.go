package command_handlers

import (
	"context"
	"errors"
	"tickets/db"
	"tickets/entities"
	"tickets/message/contracts"
)

type BookShowTicketsCommandHandler struct {
	bookingRepo contracts.BookingRepository
}

func NewBookShowTicketsCommandHandler(bookingRepo contracts.BookingRepository) BookShowTicketsCommandHandler {
	return BookShowTicketsCommandHandler{
		bookingRepo: bookingRepo,
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

	return err
}
