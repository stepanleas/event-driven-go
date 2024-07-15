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
) {
	ep.AddHandlers(
		cqrs.NewEventHandler(
			"issue_receipt_handler",
			handlers.NewIssueReceiptsHandler(receiptsService).Handle,
		),
		cqrs.NewEventHandler(
			"append_to_tracker_handler",
			handlers.NewAppendToTrackerHandler(spreadsheetsService).Handle,
		),
		cqrs.NewEventHandler(
			"tickets_to_refund_handler",
			handlers.NewTicketsToRefundHandler(spreadsheetsService).Handle,
		),
	)
}
