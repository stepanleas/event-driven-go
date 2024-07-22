package event_handlers

import (
	"context"

	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type TicketsToRefundHandler struct {
	spreadsheetsClient contracts.SpreadsheetsAPI
}

func NewTicketsToRefundHandler(spreadsheetsClient contracts.SpreadsheetsAPI) TicketsToRefundHandler {
	return TicketsToRefundHandler{spreadsheetsClient: spreadsheetsClient}
}

func (h TicketsToRefundHandler) Handle(ctx context.Context, event *entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Info("Adding ticket refund to sheet")

	return h.spreadsheetsClient.AppendRow(
		ctx,
		"tickets-to-refund",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
}
