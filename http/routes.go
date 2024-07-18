package http

import (
	"tickets/message/contracts"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	spreadsheetsAPIClient contracts.SpreadsheetsAPI,
	ticketRepo contracts.TicketRepository,
	showRepo contracts.ShowRepository,
) *echo.Echo {
	ticketCtrl := NewTicketController(eventBus, ticketRepo)
	showCtrl := NewShowController(showRepo)

	e := libHttp.NewEcho()
	e.GET("/tickets", ticketCtrl.FindAll)
	e.GET("/health", ticketCtrl.HealthCheck)
	e.POST("/tickets-status", ticketCtrl.Status)
	e.POST("/shows", showCtrl.Store)

	return e
}
