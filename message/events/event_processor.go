package events

import (
	"tickets/message/contracts"
	"tickets/message/handlers"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func AddEventProcessorHandlers(
	ep *cqrs.EventProcessor,
	eventBus *cqrs.EventBus,
	receiptsService contracts.ReceiptsService,
	spreadsheetsService contracts.SpreadsheetsAPI,
	ticketRepo contracts.TicketRepository,
	showRepo contracts.ShowRepository,
	filesAPI contracts.FilesAPI,
	deadNationAPI contracts.DeadNationApi,
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
			"RefundTicket",
			handlers.NewTicketsToRefundHandler(spreadsheetsService).Handle,
		),
		cqrs.NewEventHandler(
			"StoreTicket",
			handlers.NewStoreTicketHandler(ticketRepo).Handle,
		),
		cqrs.NewEventHandler(
			"RemoveCanceledTicket",
			handlers.NewRemoveCanceledTicketHandler(ticketRepo).Handle,
		),
		cqrs.NewEventHandler(
			"PrintTicketHandler",
			handlers.NewPrintTicketHandler(filesAPI, eventBus).Handle,
		),
		cqrs.NewEventHandler(
			"BookPlaceInDeadNation",
			handlers.NewBookingMadeHandler(deadNationAPI, showRepo).Handle,
		),
	)
}
