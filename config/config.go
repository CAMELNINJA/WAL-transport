package config

import (
	"time"

	"github.com/CAMELNINGA/cdc-postgres.git/pkg/postgres"
)

type Config struct {
	Database postgres.DatabaseCfg
	Listener Listener
	Kafka    Kafka
}

type Listener struct {
	RefreshConnection time.Duration `long:"refresh-connection" env:"REFRESH_CONNECTION" description:"Refresh connection"`
	SlotName          string        `long:"slot-name" env:"SLOT_NAME" description:"Slot name"`
}

// LoggerCfg path of the logger config.
type LoggerCfg struct {
	Caller bool   `long:"caller" env:"CALLER" description:"Caller"`
	Level  string `long:"level" env:"LEVEL" description:"Logger level"`
	Format string `long:"format" env:"FORMAT" description:"Logger format"`
}

type Kafka struct {
	Brokers []string `long:"brokers" env:"BROKERS" env-delim:"," description:"Kafka brokers"`
	Topic   string   `long:"topic" env:"TOPIC" description:"Kafka topic"`
	GroupID string   `long:"group-id" env:"GROUP_ID" description:"Kafka group id"`
}
