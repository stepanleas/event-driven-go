package http

import (
	"tickets/message/contracts"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(publisher message.Publisher, spreadsheetsAPIClient contracts.SpreadsheetsAPI) *echo.Echo {
	ctrl := NewTicketController(publisher)

	e := libHttp.NewEcho()
	e.GET("/health", ctrl.HealthCheck)
	e.POST("/tickets-status", ctrl.Status)

	return e
}
