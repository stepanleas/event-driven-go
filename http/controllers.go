package http

import (
	"fmt"
	"net/http"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

type ticketsStatusRequest struct {
	Tickets []ticketStatus `json:"tickets"`
}

type ticketStatus struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	CustomerEmail string         `json:"customer_email"`
	Price         entities.Money `json:"price"`
}

type TicketController struct {
	eventBus *cqrs.EventBus
	repo     contracts.TicketRepository
}

func NewTicketController(eventBus *cqrs.EventBus, repo contracts.TicketRepository) TicketController {
	return TicketController{eventBus: eventBus, repo: repo}
}

func (ctrl TicketController) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (ctrl TicketController) FindAll(c echo.Context) error {
	tickets, err := ctrl.repo.FindAll(c.Request().Context())
	if err != nil {
		return fmt.Errorf("failed to find tickets: %w", err)
	}

	return c.JSON(http.StatusOK, tickets)
}

func (ctrl TicketController) Status(c echo.Context) error {
	var request ticketsStatusRequest
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

			if err := ctrl.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("failed to publish TicketBookingConfirmed event: %w", err)
			}
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			if err := ctrl.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("failed to publish TicketBookingCanceled event: %w", err)
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}
