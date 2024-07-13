package handlers

import (
	"encoding/json"

	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/message"
)

const brokenMessageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

type IssueReceiptsHandler struct {
	receiptsClient contracts.ReceiptsService
}

func NewIssueReceiptsHandler(receiptsClient contracts.ReceiptsService) IssueReceiptsHandler {
	return IssueReceiptsHandler{receiptsClient: receiptsClient}
}

func (h IssueReceiptsHandler) Handle(msg *message.Message) error {
	if msg.UUID == brokenMessageID {
		return nil
	}

	if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
		return nil
	}

	var event entities.TicketBookingConfirmed
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	currency := event.Price.Currency
	if currency == "" {
		currency = "USD"
	}

	_, err := h.receiptsClient.IssueReceipt(msg.Context(), entities.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price: entities.Money{
			Amount:   event.Price.Amount,
			Currency: currency,
		},
	})

	return err
}
