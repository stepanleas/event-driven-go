package handlers

import (
	"context"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type RemoveCanceledTicketHandler struct {
	repo contracts.TicketRepository
}

func NewRemoveCanceledTicketHandler(repo contracts.TicketRepository) RemoveCanceledTicketHandler {
	return RemoveCanceledTicketHandler{repo: repo}
}

func (h RemoveCanceledTicketHandler) Handle(ctx context.Context, event *entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Info("Removing canceled ticket from database")

	return h.repo.Remove(ctx, event.TicketID)
}
