package api

import (
	"context"
	"sync"
	"tickets/entities"
	"time"
)

type ReceiptsServiceMock struct {
	mock sync.Mutex

	IssuedReceipts map[string]entities.IssueReceiptRequest
	VoidedReceipts []entities.VoidReceipt
}

func (c *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	c.mock.Lock()
	defer c.mock.Unlock()

	c.IssuedReceipts[request.TicketID] = request

	return entities.IssueReceiptResponse{
		ReceiptNumber: "mocked-receipt-number",
		IssuedAt:      time.Now(),
	}, nil
}

func (c *ReceiptsServiceMock) VoidReceipt(ctx context.Context, request entities.VoidReceipt) error {
	c.mock.Lock()
	defer c.mock.Unlock()

	c.VoidedReceipts = append(c.VoidedReceipts, request)

	return nil
}
