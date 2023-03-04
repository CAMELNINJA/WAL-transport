package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Event structure for publishing to the server.
type Event struct {
	ID        uuid.UUID      `json:"id"`
	Schema    string         `json:"schema"`
	Table     string         `json:"table"`
	Action    string         `json:"action"`
	Data      map[string]any `json:"data"`
	DataOld   map[string]any `json:"dataOld"`
	EventTime time.Time      `json:"commitTime"`
}
