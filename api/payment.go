package api

import (
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/payments"
	"golang.org/x/net/context"
)

type PaymentServiceClient struct {
	clients *clients.Clients
}

func NewPaymentServiceClient(clients *clients.Clients) PaymentServiceClient {
	if clients == nil {
		panic("NewPaymentsServiceClient: clients is nil")
	}

	return PaymentServiceClient{clients: clients}
}

func (c PaymentServiceClient) RefundPayment(ctx context.Context, refundPayment entities.PaymentRefund) error {
	resp, err := c.clients.Payments.PutRefundsWithResponse(ctx, payments.PaymentRefundRequest{
		PaymentReference: refundPayment.TicketID,
		Reason:           refundPayment.RefundReason,
		DeduplicationId:  &refundPayment.IdempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to post refund for payment %s: %w", refundPayment.TicketID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected for /payments-api/refunds status code: %d", resp.StatusCode())
	}

	return nil
}
