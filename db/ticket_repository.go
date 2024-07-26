package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	if db == nil {
		panic("db is nil")
	}

	return TicketRepository{db: db}
}

func (t TicketRepository) FindAll(ctx context.Context) ([]entities.Ticket, error) {
	var tickets []entities.Ticket

	err := t.db.SelectContext(
		ctx,
		&tickets,
		`
		SELECT 
		    ticket_id,
			price_amount AS "price.amount",
			price_currency AS "price.currency",
			customer_email 
		FROM 
		    tickets
		WHERE
			deleted_at IS NULL
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve tickets: %w", err)
	}

	return tickets, nil
}

func (t TicketRepository) Add(ctx context.Context, ticket entities.Ticket) error {
	_, err := t.db.NamedExecContext(
		ctx,
		`
		INSERT INTO
    		tickets (ticket_id, price_amount, price_currency, customer_email)
		VALUES
		    (:ticket_id, :price.amount, :price.currency, :customer_email)
		ON CONFLICT DO NOTHING
		`,
		ticket,
	)
	if err != nil {
		return fmt.Errorf("could not save ticket: %w", err)
	}

	return nil
}

func (t TicketRepository) Remove(ctx context.Context, ticketID string) error {
	res, err := t.db.ExecContext(
		ctx,
		`UPDATE tickets SET deleted_at = now() WHERE ticket_id = $1`,
		ticketID,
	)
	if err != nil {
		return fmt.Errorf("could not remove ticket: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ticket with id %s not found", ticketID)
	}

	return nil
}
