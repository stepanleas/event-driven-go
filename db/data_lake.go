package db

import (
	"context"
	"errors"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DataLake struct {
	db *sqlx.DB
}

func NewDataLake(db *sqlx.DB) DataLake {
	if db == nil {
		panic("db is nil")
	}

	return DataLake{db: db}
}

func (d DataLake) FindAll(ctx context.Context) ([]entities.DataLakeEvent, error) {
	var events []entities.DataLakeEvent
	err := d.db.SelectContext(ctx, &events, "SELECT * FROM events ORDER BY published_at ASC")
	if err != nil {
		return nil, fmt.Errorf("could not get events from data lake: %w", err)
	}

	return events, nil
}

func (d DataLake) Store(ctx context.Context, event entities.DataLakeEvent) error {
	args := map[string]interface{}{
		"event_id":      event.EventID,
		"published_at":  event.PublishedAt,
		"event_name":    event.EventName,
		"event_payload": event.EventPayload,
	}

	_, err := d.db.NamedExecContext(
		ctx,
		`
		INSERT INTO
    		events (event_id, published_at, event_name, event_payload)
		VALUES
		    (:event_id, :published_at, :event_name, :event_payload)
		`,
		args,
	)

	var postgresError *pq.Error
	if errors.As(err, &postgresError) && postgresError.Code.Name() == "unique_violation" {
		// handling re-delivery
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not store %s event in data lake: %w", event.EventID, err)
	}

	return nil
}
