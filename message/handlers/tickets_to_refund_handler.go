package handlers

import (
	"encoding/json"

	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/message"
)

type TicketsToRefundHandler struct {
	spreadsheetsClient contracts.SpreadsheetsAPI
}

func NewTicketsToRefundHandler(spreadsheetsClient contracts.SpreadsheetsAPI) TicketsToRefundHandler {
	return TicketsToRefundHandler{spreadsheetsClient: spreadsheetsClient}
}

func (h TicketsToRefundHandler) Handle(msg *message.Message) error {
	if msg.UUID == brokenMessageID {
		return nil
	}

	if msg.Metadata.Get("type") != "TicketBookingCanceled" {
		return nil
	}

	var event entities.TicketBookingCanceled
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	currency := event.Price.Currency
	if currency == "" {
		currency = "USD"
	}

	return h.spreadsheetsClient.AppendRow(
		msg.Context(),
		"tickets-to-refund",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, currency},
	)
}
