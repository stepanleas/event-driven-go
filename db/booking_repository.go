package db

import (
	"context"
	"errors"
	"fmt"
	"tickets/entities"
	"tickets/message/events"
	"tickets/message/events/outbox"

	"github.com/jmoiron/sqlx"
)

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) BookingRepository {
	if db == nil {
		panic("db is nil")
	}

	return BookingRepository{db: db}
}

func (b BookingRepository) Add(ctx context.Context, booking entities.Booking) error {
	tx, err := b.db.Beginx()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			err = errors.Join(err, rollbackErr)
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO 
		    bookings (booking_id, show_id, number_of_tickets, customer_email) 
		VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email)
		`, booking)
	if err != nil {
		return fmt.Errorf("could not add booking: %w", err)
	}

	outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
	if err != nil {
		return fmt.Errorf("could not create event bus: %w", err)
	}

	err = events.NewEventBus(outboxPublisher).Publish(ctx, entities.BookingMade{
		Header:          entities.NewEventHeader(),
		BookingID:       booking.BookingID,
		NumberOfTickets: booking.NumberOfTickets,
		CustomerEmail:   booking.CustomerEmail,
		ShowId:          booking.ShowID,
	})
	if err != nil {
		return fmt.Errorf("could not publish event: %w", err)
	}

	return nil
}
