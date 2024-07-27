package command_handlers

import (
	"context"
	"fmt"

	"tickets/entities"
	"tickets/message/command_handlers/contract"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type RefundTicketHandler struct {
	eventBus              *cqrs.EventBus
	receiptsServiceClient contract.ReceiptsService
	paymentsServiceClient contract.PaymentsService
}

func NewRefundTicketHandler(
	eventBus *cqrs.EventBus,
	receiptsServiceClient contract.ReceiptsService,
	paymentsServiceClient contract.PaymentsService,
) RefundTicketHandler {
	return RefundTicketHandler{
		eventBus:              eventBus,
		receiptsServiceClient: receiptsServiceClient,
		paymentsServiceClient: paymentsServiceClient,
	}
}

func (h RefundTicketHandler) Handle(ctx context.Context, command *entities.RefundTicket) error {
	idempotencyKey := command.Header.IdempotencyKey
	if idempotencyKey == "" {
		return fmt.Errorf("idempotency key is required")
	}

	err := h.receiptsServiceClient.VoidReceipt(ctx, entities.VoidReceipt{
		TicketID:       command.TicketID,
		Reason:         "ticket refunded",
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to void receipt: %w", err)
	}

	err = h.paymentsServiceClient.RefundPayment(ctx, entities.PaymentRefund{
		TicketID:       command.TicketID,
		RefundReason:   "ticket refunded",
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	err = h.eventBus.Publish(ctx, entities.TicketRefunded_v1{
		Header:   entities.NewEventHeader(),
		TicketID: command.TicketID,
	})
	if err != nil {
		return fmt.Errorf("failed to publish TicketRefunded event: %w", err)
	}

	return nil
}
