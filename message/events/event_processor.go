package events

import (
	"tickets/db/read_model"
	"tickets/message/contracts"
	"tickets/message/event_handlers"
	"tickets/process_manager"

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
	opsReadModel read_model.OpsBookingReadModel,
	vipBundlePM *process_manager.VipBundleProcessManager,
) {
	ep.AddHandlers(
		cqrs.NewEventHandler(
			"IssueReceipt",
			event_handlers.NewIssueReceiptsHandler(receiptsService, eventBus).Handle,
		),
		cqrs.NewEventHandler(
			"AppendToTracker",
			event_handlers.NewAppendToTrackerHandler(spreadsheetsService).Handle,
		),
		cqrs.NewEventHandler(
			"TicketRefundToSheet",
			event_handlers.NewTicketsToRefundHandler(spreadsheetsService).Handle,
		),
		cqrs.NewEventHandler(
			"StoreTicket",
			event_handlers.NewStoreTicketHandler(ticketRepo).Handle,
		),
		cqrs.NewEventHandler(
			"RemoveCanceledTicket",
			event_handlers.NewRemoveCanceledTicketHandler(ticketRepo).Handle,
		),
		cqrs.NewEventHandler(
			"PrintTicket",
			event_handlers.NewPrintTicketHandler(filesAPI, eventBus).Handle,
		),
		cqrs.NewEventHandler(
			"BookPlaceInDeadNation",
			event_handlers.NewBookingMadeHandler(deadNationAPI, showRepo).Handle,
		),
		// read model
		cqrs.NewEventHandler(
			"ops_read_model.OnBookingMade",
			opsReadModel.OnBookingMade,
		),
		cqrs.NewEventHandler(
			"ops_read_model.OnTicketReceiptIssued",
			opsReadModel.OnTicketReceiptIssued,
		),
		cqrs.NewEventHandler(
			"ops_read_model.OnTicketBookingConfirmed",
			opsReadModel.OnTicketBookingConfirmed,
		),
		cqrs.NewEventHandler(
			"ops_read_model.OnTicketPrinted",
			opsReadModel.OnTicketPrinted,
		),
		cqrs.NewEventHandler(
			"ops_read_model.OnTicketRefunded",
			opsReadModel.OnTicketRefunded,
		),
		// process manager
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnVipBundleInitialized",
			vipBundlePM.OnVipBundleInitialized,
		),
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnBookingMade",
			vipBundlePM.OnBookingMade,
		),
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnBookingFailed",
			vipBundlePM.OnBookingFailed,
		),
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnTicketBookingConfirmed",
			vipBundlePM.OnTicketBookingConfirmed,
		),
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnFlightBooked",
			vipBundlePM.OnFlightBooked,
		),
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnFlightBookingFailed",
			vipBundlePM.OnFlightBookingFailed,
		),
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnTaxiBooked",
			vipBundlePM.OnTaxiBooked,
		),
		cqrs.NewEventHandler(
			"vip_bundle_process_manager.OnTaxiBookingFailed",
			vipBundlePM.OnTaxiBookingFailed,
		),
	)
}
