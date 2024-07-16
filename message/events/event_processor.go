package events

import (
	"tickets/message/contracts"
	"tickets/message/handlers"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func AddEventProcessorHandlers(
	ep *cqrs.EventProcessor,
	receiptsService contracts.ReceiptsService,
	spreadsheetsService contracts.SpreadsheetsAPI,
	repo contracts.TicketRepository,
) {
	ep.AddHandlers(
		cqrs.NewEventHandler(
			"IssueReceipt",
			handlers.NewIssueReceiptsHandler(receiptsService).Handle,
		),
		cqrs.NewEventHandler(
			"AppendToTracker",
			handlers.NewAppendToTrackerHandler(spreadsheetsService).Handle,
		),
		cqrs.NewEventHandler(
			"StoreTicket",
			handlers.NewStoreTicketHandler(repo).Handle,
		),
		cqrs.NewEventHandler(
			"RefundTicket",
			handlers.NewTicketsToRefundHandler(spreadsheetsService).Handle,
		),
	)
}
