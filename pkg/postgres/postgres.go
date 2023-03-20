package postgres

import (
	"fmt"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type DatabaseCfg struct {
	Host     string `mapstructure:"host" valid:"required"`
	Port     uint16 `mapstructure:"port" valid:"required"`
	Name     string `mapstructure:"name" valid:"required"`
	User     string `mapstructure:"user" valid:"required"`
	Password string `mapstructure:"password" valid:"required"`
}

// initPgxConnections initialise db and replication connections.
func InitPgxConnections(cfg DatabaseCfg) (*pgx.Conn, *pgx.ReplicationConn, error) {
	pgxConf := pgx.ConnConfig{
		// TODO logger
		LogLevel: pgx.LogLevelInfo,
		Logger:   pgxLogger{},
		Host:     cfg.Host,
		Port:     cfg.Port,
		Database: cfg.Name,
		User:     cfg.User,
		Password: cfg.Password,
	}

	pgConn, err := pgx.Connect(pgxConf)
	if err != nil {
		return nil, nil, fmt.Errorf("db connection: %w", err)
	}

	rConnection, err := pgx.ReplicationConnect(pgxConf)
	if err != nil {
		return nil, nil, fmt.Errorf("replication connect: %w", err)
	}

	return pgConn, rConnection, nil
}

type pgxLogger struct{}

func (l pgxLogger) Log(level pgx.LogLevel, msg string, data map[string]any) {
	logrus.Debugln(msg)
}
