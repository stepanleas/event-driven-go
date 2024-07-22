package event_handlers

import (
	"context"
	"fmt"

	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type BookPlaceInDeadNationHandler struct {
	deadNationClient contracts.DeadNationApi
	showRepo         contracts.ShowRepository
}

func NewBookingMadeHandler(deadNationClient contracts.DeadNationApi, showRepo contracts.ShowRepository) BookPlaceInDeadNationHandler {
	return BookPlaceInDeadNationHandler{deadNationClient: deadNationClient, showRepo: showRepo}
}

func (h BookPlaceInDeadNationHandler) Handle(ctx context.Context, event *entities.BookingMade) error {
	log.FromContext(ctx).Info("Generating ticket for booking")

	show, err := h.showRepo.FindByID(ctx, event.ShowId)
	if err != nil {
		return fmt.Errorf("failed to get show: %w", err)
	}

	err = h.deadNationClient.BookInDeadNation(ctx, entities.DeadNationBooking{
		CustomerEmail:     event.CustomerEmail,
		DeadNationEventID: show.DeadNationID,
		NumberOfTickets:   event.NumberOfTickets,
		BookingID:         event.BookingID,
	})
	if err != nil {
		return fmt.Errorf("failed to book in dead nation: %w", err)
	}

	return nil
}
