package handlers

import (
	"encoding/json"
	"tickets/api"
	"tickets/events"
	"tickets/valueobject"

	"github.com/ThreeDotsLabs/watermill/message"
)

const brokenMessageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

type IssueReceiptsHandler struct {
	receiptsClient api.ReceiptsClient
}

func NewIssueReceiptsHandler(receiptsClient api.ReceiptsClient) IssueReceiptsHandler {
	return IssueReceiptsHandler{receiptsClient: receiptsClient}
}

func (h IssueReceiptsHandler) Handle(msg *message.Message) error {
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

	return h.receiptsClient.IssueReceipt(msg.Context(), api.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price: valueobject.Money{
			Amount:   event.Price.Amount,
			Currency: currency,
		},
	})
}
