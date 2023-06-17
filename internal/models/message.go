package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID         uuid.UUID      `json:"id"`
	Schema     string         `json:"schema"`
	Table      string         `json:"table"`
	Action     string         `json:"action"`
	Data       map[string]any `json:"data"`
	DataOld    map[string]any `json:"dataOld"`
	CommitTime time.Time      `json:"commitTime"`
}
