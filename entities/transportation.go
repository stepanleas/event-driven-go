package entities

import (
	"errors"

	"github.com/google/uuid"
)

var ErrNoFlightTicketsAvailable = errors.New("no flight tickets available")

type BookFlightTicketRequest struct {
	CustomerEmail  string
	FlightID       uuid.UUID
	PassengerNames []string
	ReferenceId    string
	IdempotencyKey string
}

type BookFlightTicketResponse struct {
	TicketIds []uuid.UUID `json:"ticket_ids"`
}

type BookTaxiRequest struct {
	CustomerEmail      string
	NumberOfPassengers int
	PassengerName      string
	ReferenceId        string
	IdempotencyKey     string
}

type BookTaxiResponse struct {
	TaxiBookingId uuid.UUID `json:"taxi_booking_id"`
}
