package handlers

import (
	"context"
	"fmt"

	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type IssueReceiptsHandler struct {
	receiptsClient contracts.ReceiptsService
}

func NewIssueReceiptsHandler(receiptsClient contracts.ReceiptsService) IssueReceiptsHandler {
	return IssueReceiptsHandler{receiptsClient: receiptsClient}
}

func (h IssueReceiptsHandler) Handle(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Issuing receipt")

	request := entities.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price:    event.Price,
	}

	_, err := h.receiptsClient.IssueReceipt(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	return nil
}
