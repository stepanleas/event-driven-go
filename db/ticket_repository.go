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
	_, err := t.db.ExecContext(
		ctx,
		`DELETE FROM TICKETS WHERE ticket_id = $1`,
		ticketID,
	)
	if err != nil {
		return fmt.Errorf("could not delete ticket: %w", err)
	}

	return nil
}
