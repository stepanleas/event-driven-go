package http

import (
	"errors"
	"fmt"
	"net/http"
	"tickets/db"
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

type bookTicketResponse struct {
	BookingId uuid.UUID   `json:"booking_id"`
	TicketIds []uuid.UUID `json:"ticket_ids"`
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

	bookingID := uuid.New()

	if request.NumberOfTickets < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "number of tickets must be greater than 0")
	}

	err := ctrl.repo.Add(c.Request().Context(), entities.Booking{
		BookingID:       bookingID,
		ShowID:          request.ShowID,
		NumberOfTickets: request.NumberOfTickets,
		CustomerEmail:   request.CustomerEmail,
	})
	if errors.Is(err, db.ErrNoPlacesLeft) {
		return echo.NewHTTPError(http.StatusBadRequest, "not enough seats available")
	}
	if err != nil {
		return fmt.Errorf("failed to store booking: %w", err)
	}

	return c.JSON(http.StatusCreated, bookTicketResponse{
		BookingId: bookingID,
	})
}
