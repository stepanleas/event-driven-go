package events

import (
	"tickets/valueobject"
	"time"

	"github.com/ThreeDotsLabs/watermill"
)

type EventHeader struct {
	ID          string `json:"id"`
	PublishedAt string `json:"published_at"`
}

func NewEventHeader() EventHeader {
	return EventHeader{
		ID:          watermill.NewUUID(),
		PublishedAt: time.Now().Format(time.RFC3339),
	}
}

type TicketBookingConfirmed struct {
	Header        EventHeader       `json:"header"`
	TicketID      string            `json:"ticket_id"`
	CustomerEmail string            `json:"customer_email"`
	Price         valueobject.Money `json:"price"`
}

type TicketBookingCanceled struct {
	Header        EventHeader       `json:"header"`
	TicketID      string            `json:"ticket_id"`
	CustomerEmail string            `json:"customer_email"`
	Price         valueobject.Money `json:"price"`
}
