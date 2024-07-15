package http

import (
	"tickets/message/contracts"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(eventBus *cqrs.EventBus, spreadsheetsAPIClient contracts.SpreadsheetsAPI) *echo.Echo {
	ctrl := NewTicketController(eventBus)

	e := libHttp.NewEcho()
	e.GET("/health", ctrl.HealthCheck)
	e.POST("/tickets-status", ctrl.Status)

	return e
}
