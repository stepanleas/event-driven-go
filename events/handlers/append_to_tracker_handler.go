package handlers

import (
	"encoding/json"
	"tickets/api"
	"tickets/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

type AppendToTrackerHandler struct {
	spreadsheetsClient api.SpreadsheetsClient
}

func NewAppendToTrackerHandler(spreadsheetsClient api.SpreadsheetsClient) AppendToTrackerHandler {
	return AppendToTrackerHandler{spreadsheetsClient: spreadsheetsClient}
}

func (h AppendToTrackerHandler) Handle(msg *message.Message) error {
	if msg.UUID == brokenMessageID {
		return nil
	}

	if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
		return nil
	}

	var event events.TicketBookingConfirmed
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	currency := event.Price.Currency
	if currency == "" {
		currency = "USD"
	}

	return h.spreadsheetsClient.AppendRow(
		msg.Context(),
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, currency},
	)
}
