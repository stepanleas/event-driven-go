package http

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo, pub message.Publisher) {
	ctrl := NewTicketController(pub)

	e.GET("/health", ctrl.HealthCheck)
	e.POST("/tickets-status", ctrl.Status)
}
