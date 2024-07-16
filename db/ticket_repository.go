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

func (t TicketRepository) Add(ctx context.Context, ticket entities.Ticket) error {
	_, err := t.db.NamedExecContext(
		ctx,
		`
		INSERT INTO 
    		tickets (ticket_id, price_amount, price_currency, customer_email) 
		VALUES 
		    (:ticket_id, :price.amount, :price.currency, :customer_email)`,
		ticket,
	)
	if err != nil {
		return fmt.Errorf("could not save ticket: %w", err)
	}

	return nil
}
