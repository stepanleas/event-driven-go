package http

import (
	"fmt"
	"net/http"
	"tickets/entities"
	"tickets/message/contracts"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type bookTicketRequest struct {
	ShowID          uuid.UUID `json:"show_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email"`
}

type BookingController struct {
	repo contracts.BookingRepository
}

func NewBookingController(repo contracts.BookingRepository) BookingController {
	return BookingController{repo: repo}
}

func (ctrl BookingController) Store(c echo.Context) error {
	var request bookTicketRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	booking := entities.Booking{
		BookingID:       uuid.New(),
		ShowID:          request.ShowID,
		NumberOfTickets: request.NumberOfTickets,
		CustomerEmail:   request.CustomerEmail,
	}

	if err := ctrl.repo.Add(c.Request().Context(), booking); err != nil {
		return fmt.Errorf("failed to store booking: %w", err)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"booking_id": booking.BookingID.String(),
	})
}
