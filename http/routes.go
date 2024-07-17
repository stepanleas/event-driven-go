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
	repo contracts.TicketRepository,
) *echo.Echo {
	ctrl := NewTicketController(eventBus, repo)

	e := libHttp.NewEcho()
	e.GET("/tickets", ctrl.FindAll)
	e.GET("/health", ctrl.HealthCheck)
	e.POST("/tickets-status", ctrl.Status)

	return e
}
