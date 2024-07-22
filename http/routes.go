package http

import (
	"tickets/message/contracts"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	commandBus *cqrs.CommandBus,
	spreadsheetsAPIClient contracts.SpreadsheetsAPI,
	ticketRepo contracts.TicketRepository,
	showRepo contracts.ShowRepository,
	bookingRepo contracts.BookingRepository,
) *echo.Echo {
	ticketCtrl := NewTicketController(eventBus, commandBus, ticketRepo)
	showCtrl := NewShowController(showRepo)
	bookingCtrl := NewBookingController(bookingRepo)

	e := libHttp.NewEcho()
	e.GET("/tickets", ticketCtrl.FindAll)
	e.GET("/health", ticketCtrl.HealthCheck)
	e.POST("/tickets-status", ticketCtrl.Status)
	e.PUT("/ticket-refund/:ticket_id", ticketCtrl.Refund)
	e.POST("/shows", showCtrl.Store)
	e.POST("/book-tickets", bookingCtrl.Store)

	return e
}
