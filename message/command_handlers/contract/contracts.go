package contract

import (
	"context"
	"tickets/entities"
)

type ReceiptsService interface {
	VoidReceipt(ctx context.Context, request entities.VoidReceipt) error
}

type PaymentsService interface {
	RefundPayment(ctx context.Context, request entities.PaymentRefund) error
}
