package api

import (
	"context"
	"sync"
	"tickets/entities"
	"time"
)

type ReceiptsServiceMock struct {
	mock           sync.Mutex
	IssuedReceipts []entities.IssueReceiptRequest
}

func (m *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	m.mock.Lock()
	defer m.mock.Unlock()

	m.IssuedReceipts = append(m.IssuedReceipts, request)

	return entities.IssueReceiptResponse{
		ReceiptNumber: "123",
		IssuedAt:      time.Now(),
	}, nil
}
