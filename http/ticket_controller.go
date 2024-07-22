package http

import (
	"fmt"
	"net/http"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
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
	eventBus   *cqrs.EventBus
	commandBus *cqrs.CommandBus
	repo       contracts.TicketRepository
}

func NewTicketController(eventBus *cqrs.EventBus, commandBus *cqrs.CommandBus, repo contracts.TicketRepository) TicketController {
	return TicketController{
		eventBus:   eventBus,
		commandBus: commandBus,
		repo:       repo,
	}
}

func (ctrl TicketController) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (ctrl TicketController) Refund(c echo.Context) error {
	ticketID := c.Param("ticket_id")

	command := entities.RefundTicket{
		Header:   entities.NewEventHeaderWithIdempotencyKey(uuid.NewString()),
		TicketID: ticketID,
	}

	if err := ctrl.commandBus.Send(c.Request().Context(), command); err != nil {
		return fmt.Errorf("failed to send RefundTicket command: %w", err)
	}

	return c.NoContent(http.StatusAccepted)
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

	idempotencyKey := c.Request().Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Idempotency-Key header is required")
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			event := entities.TicketBookingConfirmed{
				Header:        entities.NewEventHeaderWithIdempotencyKey(idempotencyKey + ticket.TicketID),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			if err := ctrl.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("failed to publish TicketBookingConfirmed event: %w", err)
			}
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewEventHeaderWithIdempotencyKey(idempotencyKey + ticket.TicketID),
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
