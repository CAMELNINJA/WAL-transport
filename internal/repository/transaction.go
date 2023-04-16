package repository

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tx struct {
	tx pgx.Tx
	mu sync.Mutex
}

func NewTx(ctx context.Context, p *pgxpool.Pool) (*Tx, error) {
	tx, err := p.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

func (ct *Tx) Rollback(ctx context.Context) error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	return ct.tx.Rollback(ctx)
}

func (ct *Tx) Commit(ctx context.Context) error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	return ct.tx.Commit(ctx)
}

func (ct *Tx) QueryRow(ctx context.Context, sql string, args ...interface{}) error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	_, err := ct.tx.Exec(ctx, sql, args...)

	return err
}

// query rebinds the query and executes it.
func (ct *Tx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	return ct.tx.Query(ctx, sql, args...)
}
