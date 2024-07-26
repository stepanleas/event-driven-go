package entities

import (
	"time"

	"github.com/google/uuid"
)

type DataLakeEvent struct {
	EventID      uuid.UUID `json:"event_id"`
	PublishedAt  time.Time `json:"published_at"`
	EventName    string    `json:"event_name"`
	EventPayload string    `json:"event_payload"`
}
