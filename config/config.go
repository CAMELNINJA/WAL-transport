package config

import (
	"time"

	"github.com/CAMELNINGA/WAL-transport.git/pkg/postgres"
)

const (
	FilterType  = "filter"
	ReplaseType = "replase"
)

type Config struct {
	Database postgres.DatabaseCfg `json:"database"`
	Listener Listener             `json:"listener"`
	Kafka    Kafka                `json:"kafka"`
	Sanitize []Sanitize           `json:"sanitize"`
}

type Listener struct {
	RefreshConnection time.Duration `json:"refresh_connection"`
	SlotName          string        `json:"slot_name"`
}

// LoggerCfg path of the logger config.
type LoggerCfg struct {
	Caller bool   `long:"caller" env:"CALLER" description:"Caller"`
	Level  string `long:"level" env:"LEVEL" description:"Logger level" default:"info"`
	Format string `long:"format" env:"FORMAT" description:"Logger format"`
}

type Kafka struct {
	Brokers []string `json:"brokers" long:"brokers" env:"BROKERS" env-delim:"," description:"Kafka brokers"`
	Topic   string   `json:"topic" long:"topic" env:"TOPIC" description:"Kafka topic"`
	GroupID string   `json:"group_id" long:"group-id" env:"GROUP_ID" description:"Kafka group id"`
}

type Sanitize struct {
	Type     string            `json:"type" long:"type" env:"TYPE" description:"Sanitize type"`
	Table    string            `json:"table" long:"table" env:"TABLE" description:"Table name"`
	OldTable string            `json:"old_table" long:"old-table" env:"OLD_TABLE" description:"Old table name"`
	Schema   map[string]string `json:"schema" long:"schema" env:"SCHEMA" description:"Schema name"`
	Columns  map[string]string `json:"columns" long:"columns" env:"COLUMNS" description:"Columns name"`
}
