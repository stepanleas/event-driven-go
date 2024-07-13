package events

import (
	"tickets/api"
	"tickets/events/handlers"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Router struct {
	router *message.Router
}

func (r Router) GetRouter() *message.Router {
	return r.router
}

func NewRouter(logger watermill.LoggerAdapter) Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	return Router{router: router}
}

func (r Router) AddMiddleware(m ...message.HandlerMiddleware) {
	r.router.AddMiddleware()
}

func (r Router) AddIssueReceiptHandler(sub message.Subscriber, receiptsClient api.ReceiptsClient) {
	r.router.AddNoPublisherHandler(
		"issue_receipt_handler",
		"TicketBookingConfirmed",
		sub,
		handlers.NewIssueReceiptsHandler(receiptsClient).Handle,
	)
}

func (r Router) AddAppendToTrackerHandler(sub message.Subscriber, spreadsheetsClient api.SpreadsheetsClient) {
	r.router.AddNoPublisherHandler(
		"append_to_tracker_handler",
		"TicketBookingConfirmed",
		sub,
		handlers.NewAppendToTrackerHandler(spreadsheetsClient).Handle,
	)
}

func (r Router) AddTicketsToRefundHandler(sub message.Subscriber, spreadsheetsClient api.SpreadsheetsClient) {
	r.router.AddNoPublisherHandler(
		"tickets_to_refund_handler",
		"TicketBookingCanceled",
		sub,
		handlers.NewTicketsToRefundHandler(spreadsheetsClient).Handle,
	)
}
