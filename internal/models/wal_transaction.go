package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	error_walListner "github.com/CAMELNINGA/WAL-transport.git/pkg/error_walListner"
)

// ActionKind kind of action on WAL message.
type ActionKind string

// kind of WAL message.
const (
	ActionKindInsert ActionKind = "INSERT"
	ActionKindUpdate ActionKind = "UPDATE"
	ActionKindDelete ActionKind = "DELETE"
)

// PostgreSQL OIDs
// https://github.com/postgres/postgres/blob/master/src/include/catalog/pg_type.dat
const (
	Int2OID = 21
	Int4OID = 23
	Int8OID = 20

	TextOID    = 25
	VarcharOID = 1043

	TimestampOID   = 1114
	TimestamptzOID = 1184
	DateOID        = 1082
	TimeOID        = 1083

	JSONBOID = 3802
	UUIDOID  = 2950
	BoolOID  = 16
)

// WalTransaction transaction specified WAL message.
type WalTransaction struct {
	LSN           int64
	BeginTime     *time.Time
	CommitTime    *time.Time
	RelationStore map[int32]RelationData
	Actions       []*ActionData
}

func (wt *WalTransaction) String() string {
	if wt.CommitTime != nil {
		return fmt.Sprintf("CommitTime %v Actions %v ", wt.CommitTime, wt.Actions)
	}
	return fmt.Sprintf("BeginTime %v Actions %v ", wt.BeginTime, wt.Actions)
}

// NewWalTransaction create and initialize new WAL transaction.
func NewWalTransaction() *WalTransaction {
	return &WalTransaction{
		RelationStore: make(map[int32]RelationData),
	}
}

func (k ActionKind) string() string {
	return string(k)
}

// RelationData kind of WAL message data.
type RelationData struct {
	Schema  string
	Table   string
	Columns []Column
}

// ActionData kind of WAL message data.
type ActionData struct {
	Schema     string
	Table      string
	Kind       ActionKind
	OldColumns []Column
	NewColumns []Column
}

func (a ActionData) String() string {
	return fmt.Sprintf("\n Schema %s Table %s  Kind %s \n NewColumns  %v \n OldColums %v", a.Schema, a.Table, a.Kind, a.NewColumns, a.OldColumns)
}

// Column of the table with which changes occur.
type Column struct {
	Name      string
	Value     any
	ValueType int
	IsKey     bool
}

func (c Column) String() string {
	return fmt.Sprintf(" Name %s Value %v ValueType %d IsKey %t", c.Name, c.Value, c.ValueType, c.IsKey)
}

// AssertValue converts bytes to a specific type depending
// on the type of this data in the database table.
func (c *Column) AssertValue(src []byte) {
	var (
		val any
		err error
	)

	if src == nil {
		c.Value = nil
		return
	}

	strSrc := string(src)

	const (
		timestampLayout       = "2006-01-02 15:04:05"
		timestampWithTZLayout = "2006-01-02 15:04:05.999999999-07"
	)

	switch c.ValueType {
	case BoolOID:
		val, err = strconv.ParseBool(strSrc)
	case Int2OID, Int4OID:
		val, err = strconv.Atoi(strSrc)
	case Int8OID:
		val, err = strconv.ParseInt(strSrc, 10, 64)
	case TextOID, VarcharOID:
		val = strSrc
	case TimestampOID:
		val, err = time.Parse(timestampLayout, strSrc)
	case TimestamptzOID:
		val, err = time.ParseInLocation(timestampWithTZLayout, strSrc, time.UTC)
	case DateOID, TimeOID:
		val = strSrc
	case UUIDOID:
		val, err = uuid.Parse(strSrc)
	case JSONBOID:
		var m any
		if src[0] == '[' {
			m = make([]any, 0)
		} else {
			m = make(map[string]any)
		}
		err = json.Unmarshal(src, &m)
		val = m
	default:
		logrus.WithFields(logrus.Fields{"pgtype": c.ValueType, "column_name": c.Name}).Warnln("unknown oid type")
		val = strSrc
	}

	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"pgtype": c.ValueType, "column_name": c.Name}).
			Errorln("column data parse error")
	}

	c.Value = val
}

// Clear transaction data.
func (w *WalTransaction) Clear() {
	w.CommitTime = nil
	w.BeginTime = nil
	w.Actions = nil
}

// CreateActionData create action  from WAL message data.
func (w *WalTransaction) CreateActionData(relationID int32, oldRows []TupleData, newRows []TupleData, kind ActionKind) (a ActionData, err error) {
	rel, ok := w.RelationStore[relationID]
	if !ok {
		return a, error_walListner.ErrRelationNotFound
	}

	a = ActionData{
		Schema: rel.Schema,
		Table:  rel.Table,
		Kind:   kind,
	}

	var oldColumns []Column

	for num, row := range oldRows {
		column := Column{
			Name:      rel.Columns[num].Name,
			ValueType: rel.Columns[num].ValueType,
			IsKey:     rel.Columns[num].IsKey,
		}
		column.AssertValue(row.Value)
		oldColumns = append(oldColumns, column)
	}

	a.OldColumns = oldColumns

	var newColumns []Column
	for num, row := range newRows {
		column := Column{
			Name:      rel.Columns[num].Name,
			ValueType: rel.Columns[num].ValueType,
			IsKey:     rel.Columns[num].IsKey,
		}
		column.AssertValue(row.Value)
		newColumns = append(newColumns, column)
	}
	a.NewColumns = newColumns
	fmt.Println(a)
	return a, nil
}

// CreateMessges message by table,
// action and create messages for each value.
func (w *WalTransaction) CreateMessges() []Message {
	var messages []Message

	for _, item := range w.Actions {
		dataOld := make(map[string]any)
		for _, val := range item.OldColumns {
			dataOld[val.Name] = val.Value
		}

		data := make(map[string]any)
		for _, val := range item.NewColumns {
			data[val.Name] = val.Value
		}

		x := uuid.New()
		message := Message{
			ID:         x,
			Schema:     item.Schema,
			Table:      item.Table,
			Action:     item.Kind.string(),
			DataOld:    dataOld,
			Data:       data,
			CommitTime: *w.CommitTime,
		}
		messages = append(messages, message)
		// filterSkippedEvents.With(prometheus.Labels{"table": item.Table}).Inc()

		// logrus.WithFields(
		// 	logrus.Fields{
		// 		"schema": item.Schema,
		// 		"table":  item.Table,
		// 		"action": item.Kind,
		// 	}).
		// 	Infoln("wal-message was skipped by filter")
	}

	return messages
}

// inArray checks whether the value is in an array.
func inArray(arr []string, value string) bool {
	for _, v := range arr {
		if strings.EqualFold(v, value) {
			return true
		}
	}

	return false
}
