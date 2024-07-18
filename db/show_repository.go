package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShowRepository struct {
	db *sqlx.DB
}

func NewShowRepository(db *sqlx.DB) ShowRepository {
	if db == nil {
		panic("db is nil")
	}

	return ShowRepository{db: db}
}

func (s ShowRepository) Add(ctx context.Context, show entities.Show) error {
	_, err := s.db.NamedExecContext(
		ctx,
		`
		INSERT INTO
    		shows (show_id, dead_nation_id, number_of_tickets, start_time, title, venue)
		VALUES
		    (:show_id, :dead_nation_id, :number_of_tickets, :start_time, :title, :venue)
		ON CONFLICT DO NOTHING
		`,
		show,
	)
	if err != nil {
		return fmt.Errorf("could not save show: %w", err)
	}

	return nil
}

func (s ShowRepository) FindAll(ctx context.Context) ([]entities.Show, error) {
	var shows []entities.Show
	err := s.db.SelectContext(ctx, &shows, `
		SELECT 
		    * 
		FROM 
		    shows
	`)
	if err != nil {
		return nil, fmt.Errorf("could not get shows: %w", err)
	}

	return shows, nil
}

func (s ShowRepository) FindByID(ctx context.Context, showID uuid.UUID) (entities.Show, error) {
	var show entities.Show
	err := s.db.GetContext(ctx, &show, `
		SELECT 
		    * 
		FROM 
		    shows
		WHERE
		    show_id = $1
	`, showID)
	if err != nil {
		return entities.Show{}, fmt.Errorf("could not get show: %w", err)
	}

	return show, nil
}
