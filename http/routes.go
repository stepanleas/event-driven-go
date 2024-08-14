package http

import (
	"tickets/db/read_model"
	"tickets/message/contracts"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

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
	opsReadModel read_model.OpsBookingReadModel,
) *echo.Echo {
	ticketCtrl := NewTicketController(eventBus, commandBus, ticketRepo)
	showCtrl := NewShowController(showRepo)
	bookingCtrl := NewBookingController(bookingRepo)
	opsBookingCtrl := NewOpsBookingController(opsReadModel)

	e := libHttp.NewEcho()

	e.Use(otelecho.Middleware("tickets"))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/health", ticketCtrl.HealthCheck)

	e.GET("/tickets", ticketCtrl.FindAll)
	e.POST("/tickets-status", ticketCtrl.Status)
	e.POST("/book-tickets", bookingCtrl.Store)
	e.PUT("/ticket-refund/:ticket_id", ticketCtrl.Refund)

	e.GET("/shows", showCtrl.FindAll)
	e.POST("/shows", showCtrl.Store)

	e.GET("/ops/bookings", opsBookingCtrl.FindAll)
	e.GET("/ops/bookings/:id", opsBookingCtrl.FindByID)

	return e
}
