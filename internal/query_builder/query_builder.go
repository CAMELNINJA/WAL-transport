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
	query := sq.Update(tx.Table)
	for _, v := range tx.NewColumns {
		query = query.Set(v.Name, v.Value)
	}
	for _, v := range tx.OldColumns {
		query = query.Where(sq.Eq{v.Name: v.Value})
	}

	return query.ToSql()
}

func prepareQueryDelete(tx *models.ActionData) (string, []interface{}, error) {
	query := sq.Delete(tx.Table)
	if len(tx.OldColumns) == 0 {
		return "", []interface{}{}, fmt.Errorf("not build")
	}
	for _, v := range tx.OldColumns {
		query = query.Where(sq.Eq{v.Name: v.Value})
	}
	return query.ToSql()
}
