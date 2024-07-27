package event_handlers

import (
	"context"
	"fmt"

	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type IssueReceiptsHandler struct {
	receiptsClient contracts.ReceiptsService
	eventBus       *cqrs.EventBus
}

func NewIssueReceiptsHandler(receiptsClient contracts.ReceiptsService, eventBus *cqrs.EventBus) IssueReceiptsHandler {
	return IssueReceiptsHandler{
		receiptsClient: receiptsClient,
		eventBus:       eventBus,
	}
}

func (h IssueReceiptsHandler) Handle(ctx context.Context, event *entities.TicketBookingConfirmed_v1) error {
	log.FromContext(ctx).Info("Issuing receipt")

	request := entities.IssueReceiptRequest{
		IdempotencyKey: event.Header.IdempotencyKey,
		TicketID:       event.TicketID,
		Price:          event.Price,
	}

	resp, err := h.receiptsClient.IssueReceipt(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	return h.eventBus.Publish(ctx, entities.TicketReceiptIssued_v1{
		Header:        entities.NewEventHeaderWithIdempotencyKey(event.Header.IdempotencyKey),
		TicketID:      event.TicketID,
		ReceiptNumber: resp.ReceiptNumber,
		IssuedAt:      resp.IssuedAt,
	})
}
