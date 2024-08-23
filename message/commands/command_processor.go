package commands

import (
	"tickets/message/command_handlers"
	"tickets/message/command_handlers/contract"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func AddCommandProcessorHandlers(
	cp *cqrs.CommandProcessor,
	eventBus *cqrs.EventBus,
	bookingRepo contracts.BookingRepository,
	receiptsServiceClient contract.ReceiptsService,
	paymentsServiceClient contract.PaymentsService,
) {
	cp.AddHandlers(
		cqrs.NewCommandHandler(
			"TicketRefund",
			command_handlers.NewRefundTicketHandler(eventBus, receiptsServiceClient, paymentsServiceClient).Handle,
		),
		cqrs.NewCommandHandler(
			"BookShowTickets",
			command_handlers.NewBookShowTicketsCommandHandler(bookingRepo).Handle,
		),
	)
}
