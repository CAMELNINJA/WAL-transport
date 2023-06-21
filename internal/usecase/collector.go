package usecase

import (
	"context"
	"fmt"

	"github.com/CAMELNINGA/WAL-transport.git/internal/models"
	repo "github.com/CAMELNINGA/WAL-transport.git/internal/repository"
	"github.com/CAMELNINGA/WAL-transport.git/internal/sanitize"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
)

type querybuilder interface {
	QueryBuilder(tx *models.ActionData) (string, []interface{}, error)
}

type transaction interface {
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) error
}

type collector struct {
	querybuilder querybuilder
	connections  *pgxpool.Pool
	sanitize     sanitize.Handler
}

func NewCollector(qb querybuilder,
	connections *pgxpool.Pool,
	sanitizer sanitize.Handler,
) *collector {
	return &collector{
		querybuilder: qb,
		connections:  connections,
		sanitize:     sanitizer,
	}
}

func checkTransaction(posTX transaction, err error) error {

	if err != nil {
		fmt.Println(err)
		return posTX.Rollback(context.Background())
	}
	return posTX.Commit(context.Background())

}

func (c *collector) SaveData(ctx context.Context, message models.Message) error {
	tx := message.ToWalTransaction()
	posTX, err := repo.NewTx(ctx, c.connections)
	if err != nil {
		return err
	}
	for _, v := range tx.Actions {
		v1 := c.sanitize.Handle(v)
		if v1 == nil {

			continue
		}
		sql, args, err := c.querybuilder.QueryBuilder(v1)
		if err != nil {
			return checkTransaction(posTX, err)
		}
		sql = sqlx.Rebind(sqlx.DOLLAR, sql)
		err = posTX.QueryRow(ctx, sql, args...)
		if err != nil {
			return checkTransaction(posTX, err)
		}
	}
	return checkTransaction(posTX, nil)
}
