package contracts

import (
	"context"
	"tickets/entities"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type TicketRepository interface {
	Add(ctx context.Context, ticket entities.Ticket) error
}
