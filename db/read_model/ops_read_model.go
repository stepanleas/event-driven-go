package read_model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"tickets/entities"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/jmoiron/sqlx"
)

type OpsBookingReadModel struct {
	db *sqlx.DB
}

func NewOpsBookingReadModel(db *sqlx.DB) OpsBookingReadModel {
	return OpsBookingReadModel{db: db}
}

func (r OpsBookingReadModel) AllReservations(receiptIssueDateFilter string) ([]entities.OpsBooking, error) {
	query := "SELECT payload FROM read_model_ops_bookings"
	var queryArgs []any

	if receiptIssueDateFilter != "" {
		query += `
			WHERE booking_id IN (
				SELECT booking_id FROM (
					SELECT booking_id, 
						DATE(jsonb_path_query(payload, '$.tickets.*.receipt_issued_at')::text) as receipt_issued_at 
					FROM 
						read_model_ops_bookings
				) bookings_within_date 
				WHERE receipt_issued_at = $1
			)
		`
		queryArgs = append(queryArgs, receiptIssueDateFilter)
	}

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.OpsBooking
	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			return nil, err
		}

		var reservation entities.OpsBooking
		if err := json.Unmarshal(payload, &reservation); err != nil {
			return nil, err
		}

		result = append(result, reservation)
	}

	return result, nil
}

func (r OpsBookingReadModel) BookingReadModel(ctx context.Context, bookingID string) (entities.OpsBooking, error) {
	return r.findReadModelByBookingID(ctx, bookingID, r.db)
}

func (r OpsBookingReadModel) OnBookingMade(ctx context.Context, event *entities.BookingMade) error {
	err := r.createReadModel(ctx, entities.OpsBooking{
		BookingID:  event.BookingID,
		BookedAt:   event.Header.PublishedAt,
		Tickets:    nil,
		LastUpdate: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("could not create read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) OnTicketReceiptIssued(ctx context.Context, event *entities.TicketReceiptIssued) error {
	err := r.updateTicketInBookingReadModel(
		ctx,
		event.TicketID,
		func(ticket entities.OpsTicket) (entities.OpsTicket, error) {
			ticket.ReceiptNumber = event.ReceiptNumber
			ticket.ReceiptIssuedAt = event.IssuedAt

			return ticket, nil
		})
	if err != nil {
		return fmt.Errorf("could not update ticket in read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) OnTicketBookingConfirmed(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	return r.updateBookingReadModel(
		ctx,
		event.BookingID,
		func(rm entities.OpsBooking) (entities.OpsBooking, error) {
			ticket, ok := rm.Tickets[event.TicketID]
			if !ok {
				// we use the zero-value of OpsTicket
				log.FromContext(ctx).
					WithField("ticket_id", event.TicketID).
					Debug("Creating ticket read model for ticket %s")
			}

			ticket.PriceAmount = event.Price.Amount
			ticket.PriceCurrency = event.Price.Currency
			ticket.CustomerEmail = event.CustomerEmail
			ticket.Status = "confirmed"

			rm.Tickets[event.TicketID] = ticket

			return rm, nil
		},
	)
}

func (r OpsBookingReadModel) OnTicketPrinted(ctx context.Context, event *entities.TicketPrinted) error {
	err := r.updateTicketInBookingReadModel(
		ctx,
		event.TicketID,
		func(ticket entities.OpsTicket) (entities.OpsTicket, error) {
			ticket.PrintedAt = event.Header.PublishedAt
			ticket.PrintedFileName = event.FileName

			return ticket, nil
		})
	if err != nil {
		return fmt.Errorf("could not update ticket in read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) OnTicketRefunded(ctx context.Context, event *entities.TicketRefunded) error {
	err := r.updateTicketInBookingReadModel(
		ctx,
		event.TicketID,
		func(ticket entities.OpsTicket) (entities.OpsTicket, error) {
			ticket.Status = "refunded"

			return ticket, nil
		})
	if err != nil {
		return fmt.Errorf("could not update ticket in read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) createReadModel(ctx context.Context, booking entities.OpsBooking) error {
	payload, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO 
		    read_model_ops_bookings (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT (booking_id) DO NOTHING; -- read model may be already updated by another event - we don't want to override
	`

	_, err = r.db.ExecContext(ctx, query, payload, booking.BookingID)
	if err != nil {
		return fmt.Errorf("could not create read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) updateBookingReadModel(
	ctx context.Context,
	bookingID string,
	updateFunc func(ticket entities.OpsBooking) (entities.OpsBooking, error),
) (err error) {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			rm, err := r.findReadModelByBookingID(ctx, bookingID, tx)
			if err == sql.ErrNoRows {
				// events arrived out of order - it should spin until the read model is created
				return fmt.Errorf("read model for booking %s not exist yet", bookingID)
			} else if err != nil {
				return fmt.Errorf("could not find booking read model: %w", err)
			}

			updatedRm, err := updateFunc(rm)
			if err != nil {
				return err
			}

			return r.updateReadModel(ctx, tx, updatedRm)
		},
	)
}

func (r OpsBookingReadModel) updateTicketInBookingReadModel(
	ctx context.Context,
	ticketID string,
	updateFunc func(ticket entities.OpsTicket) (entities.OpsTicket, error),
) (err error) {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			rm, err := r.findReadModelByTicketID(ctx, ticketID, tx)
			if err == sql.ErrNoRows {
				// events arrived out of order - it should spin until the read model is created
				return fmt.Errorf("read model for ticket %s not exist yet", ticketID)
			} else if err != nil {
				return fmt.Errorf("could not find ticket read model: %w", err)
			}

			ticket := rm.Tickets[ticketID]

			updatedRm, err := updateFunc(ticket)
			if err != nil {
				return err
			}

			rm.Tickets[ticketID] = updatedRm

			return r.updateReadModel(ctx, tx, rm)
		},
	)
}

func (r OpsBookingReadModel) updateReadModel(ctx context.Context, tx *sqlx.Tx, rm entities.OpsBooking) error {
	rm.LastUpdate = time.Now()

	payload, err := json.Marshal(rm)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO 
			read_model_ops_bookings (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT (booking_id) DO UPDATE SET payload = excluded.payload;
	`

	_, err = tx.ExecContext(ctx, query, payload, rm.BookingID)
	if err != nil {
		return fmt.Errorf("could not update read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) findReadModelByTicketID(ctx context.Context, ticketID string, db dbExecutor) (entities.OpsBooking, error) {
	var payload []byte

	query := "SELECT payload FROM read_model_ops_bookings WHERE payload::jsonb -> 'tickets' ? $1"

	err := db.QueryRowContext(ctx, query, ticketID).Scan(&payload)
	if err != nil {
		return entities.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r OpsBookingReadModel) findReadModelByBookingID(ctx context.Context, bookingID string, db dbExecutor) (entities.OpsBooking, error) {
	var payload []byte

	query := "SELECT payload FROM read_model_ops_bookings WHERE booking_id = $1"

	err := db.QueryRowContext(ctx, query, bookingID).Scan(&payload)
	if err != nil {
		return entities.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r OpsBookingReadModel) unmarshalReadModelFromDB(payload []byte) (entities.OpsBooking, error) {
	var dbReadModel entities.OpsBooking
	if err := json.Unmarshal(payload, &dbReadModel); err != nil {
		return entities.OpsBooking{}, err
	}

	if dbReadModel.Tickets == nil {
		dbReadModel.Tickets = map[string]entities.OpsTicket{}
	}

	return dbReadModel, nil
}

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
