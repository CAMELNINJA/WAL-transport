package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// initPgxConnections initialise db connections.
func InitMasterConnection(ctx context.Context, cfg DatabaseCfg) (*pgxpool.Pool, error) {

	pgConn, err := pgxpool.New(ctx, ConnectString(cfg))
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}

	return pgConn, nil
}

func ConnectString(cfg DatabaseCfg) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password)
}
