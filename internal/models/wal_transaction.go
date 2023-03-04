package models

import (
	"time"
)

// ActionKind kind of action on WAL message.
type ActionKind string

// kind of WAL message.
const (
	ActionKindInsert ActionKind = "INSERT"
	ActionKindUpdate ActionKind = "UPDATE"
	ActionKindDelete ActionKind = "DELETE"
)

// WalTransaction transaction specified WAL message.
type WalTransaction struct {
	LSN           int64
	BeginTime     *time.Time
	CommitTime    *time.Time
	RelationStore map[int32]RelationData
	Actions       []ActionData
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

// Column of the table with which changes occur.
type Column struct {
	name      string
	value     any
	valueType int
	isKey     bool
}
