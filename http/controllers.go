package http

import (
	"encoding/json"
	"net/http"
	"tickets/events/entities"
	"tickets/valueobject"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatus `json:"tickets"`
}

type TicketStatus struct {
	TicketID      string            `json:"ticket_id"`
	Status        string            `json:"status"`
	CustomerEmail string            `json:"customer_email"`
	Price         valueobject.Money `json:"price"`
}

type TicketController struct {
	pub message.Publisher
}

func NewTicketController(pub message.Publisher) TicketController {
	return TicketController{pub: pub}
}

func (ctrl TicketController) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (ctrl TicketController) Status(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			event := entities.TicketBookingConfirmed{
				Header:        entities.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), []byte(payload))
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-Id"))
			msg.Metadata.Set("type", "TicketBookingConfirmed")

			if err := ctrl.pub.Publish("TicketBookingConfirmed", msg); err != nil {
				return err
			}
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), []byte(payload))
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-Id"))
			msg.Metadata.Set("type", "TicketBookingCanceled")

			if err := ctrl.pub.Publish("TicketBookingCanceled", msg); err != nil {
				return err
			}
		}

	}

	return c.NoContent(http.StatusOK)
}
