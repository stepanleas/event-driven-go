package events

import (
	"tickets/message/contracts"
	"tickets/message/event_handlers"

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
			"IssueReceiptHandler",
			event_handlers.NewIssueReceiptsHandler(receiptsService).Handle,
		),
		cqrs.NewEventHandler(
			"AppendToTrackerHandler",
			event_handlers.NewAppendToTrackerHandler(spreadsheetsService).Handle,
		),
		cqrs.NewEventHandler(
			"RefundTicketHandler",
			event_handlers.NewTicketsToRefundHandler(spreadsheetsService).Handle,
		),
		cqrs.NewEventHandler(
			"StoreTicketHandler",
			event_handlers.NewStoreTicketHandler(ticketRepo).Handle,
		),
		cqrs.NewEventHandler(
			"RemoveCanceledTicketHandler",
			event_handlers.NewRemoveCanceledTicketHandler(ticketRepo).Handle,
		),
		cqrs.NewEventHandler(
			"PrintTicketHandler",
			event_handlers.NewPrintTicketHandler(filesAPI, eventBus).Handle,
		),
		cqrs.NewEventHandler(
			"BookPlaceInDeadNationHandler",
			event_handlers.NewBookingMadeHandler(deadNationAPI, showRepo).Handle,
		),
	)
}
