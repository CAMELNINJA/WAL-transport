package querybuilder

import (
	"fmt"

	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type qb struct {
	loger *logrus.Entry
}

func NewQueryBuilder(loger *logrus.Entry) *qb {
	return &qb{
		loger: loger,
	}
}

func (qb *qb) QueryBuilder(tx *models.ActionData) (string, []interface{}, error) {
	switch tx.Kind {
	case models.ActionKindInsert:
		return prepareQueryInsert(tx)
	case models.ActionKindUpdate:
		return prepareQueryUpdate(tx)
	case models.ActionKindDelete:
		return prepareQueryDelete(tx)
	default:
		return "", []interface{}{}, fmt.Errorf("not build")
	}
}

func prepareQueryInsert(tx *models.ActionData) (string, []interface{}, error) {
	var values []interface{}
	var columns []string
	for _, v := range tx.NewColumns {
		columns = append(columns, v.Name)
		values = append(values, v.Value)
	}
	return sq.Insert(tx.Table).Columns(columns...).Values(values...).ToSql()
}

func prepareQueryUpdate(tx *models.ActionData) (string, []interface{}, error) {
	var values []interface{}
	var columns []string
	for _, v := range tx.NewColumns {
		columns = append(columns, v.Name)
		values = append(values, v.Value)
	}
	return sq.Update(tx.Table).SetMap(sq.Eq{"id": 1}).ToSql()
}

func prepareQueryDelete(tx *models.ActionData) (string, []interface{}, error) {
	return sq.Delete(tx.Table).Where(sq.Eq{"id": 1}).ToSql()
}
