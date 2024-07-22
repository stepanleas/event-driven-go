package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"tickets/entities"
	"tickets/message/events"
	"tickets/message/events/outbox"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
	tx, err := b.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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

	availableSeats := 0
	err = tx.GetContext(ctx, &availableSeats, `
		SELECT
		    number_of_tickets AS available_seats
		FROM
		    shows
		WHERE
		    show_id = $1
	`, booking.ShowID)
	if err != nil {
		return fmt.Errorf("could not get available seats: %w", err)
	}

	alreadyBookedSeats := 0
	err = tx.GetContext(ctx, &alreadyBookedSeats, `
		SELECT
		    coalesce(SUM(number_of_tickets), 0) AS already_booked_seats
		FROM
		    bookings
		WHERE
		    show_id = $1
	`, booking.ShowID)
	if err != nil {
		return fmt.Errorf("could not get already booked seats: %w", err)
	}

	if availableSeats-alreadyBookedSeats < booking.NumberOfTickets {
		// this is usually a bad idea, learn more here: https://threedots.tech/post/introducing-clean-architecture/
		return echo.NewHTTPError(http.StatusBadRequest, "not enough seats available")
	}

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
