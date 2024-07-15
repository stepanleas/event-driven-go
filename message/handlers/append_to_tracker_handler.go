package handlers

import (
	"context"

	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type AppendToTrackerHandler struct {
	spreadsheetsClient contracts.SpreadsheetsAPI
}

func NewAppendToTrackerHandler(spreadsheetsClient contracts.SpreadsheetsAPI) AppendToTrackerHandler {
	return AppendToTrackerHandler{spreadsheetsClient: spreadsheetsClient}
}

func (h AppendToTrackerHandler) Handle(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Generating ticket for booking")

	return h.spreadsheetsClient.AppendRow(
		ctx,
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
}
