package handlers

import (
	"context"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type StoreTicketHandler struct {
	repo contracts.TicketRepository
}

func NewStoreTicketHandler(repo contracts.TicketRepository) StoreTicketHandler {
	return StoreTicketHandler{repo: repo}
}

func (h StoreTicketHandler) Handle(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Storing ticket in database")

	ticket := entities.Ticket{
		TicketID:      event.TicketID,
		Price:         event.Price,
		CustomerEmail: event.CustomerEmail,
	}

	return h.repo.Add(ctx, ticket)
}
