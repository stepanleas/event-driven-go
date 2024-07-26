package entities

import (
	"time"

	"github.com/google/uuid"
)

type OpsBooking struct {
	BookingID uuid.UUID `json:"booking_id"`
	BookedAt  time.Time `json:"booked_at"`

	Tickets map[string]OpsTicket `json:"tickets"`

	LastUpdate time.Time `json:"last_update"`
}

type OpsTicket struct {
	PriceAmount   string `json:"price_amount"`
	PriceCurrency string `json:"price_currency"`
	CustomerEmail string `json:"customer_email"`

	ConfirmedAt time.Time `json:"confirmed_at"`
	RefundedAt  time.Time `json:"refunded_at"`

	PrintedAt       time.Time `json:"printed_at"`
	PrintedFileName string    `json:"printed_file_name"`

	ReceiptIssuedAt time.Time `json:"receipt_issued_at"`
	ReceiptNumber   string    `json:"receipt_number"`
}
