package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID      `json:"id"`
	Schema     string         `json:"schema"`
	Table      string         `json:"table"`
	Action     string         `json:"action"`
	Data       map[string]any `json:"data"`
	DataOld    map[string]any `json:"dataOld"`
	CommitTime time.Time      `json:"commitTime"`
}

func (m Message) ToWalTransaction() *WalTransaction {
	return &WalTransaction{
		Actions: []ActionData{m.ToActionData()},
	}
}

func (m Message) ToActionData() ActionData {
	return ActionData{
		Schema:     m.Schema,
		Table:      m.Table,
		Kind:       ActionKind(m.Action),
		OldColumns: m.ToColumns(m.DataOld),
		NewColumns: m.ToColumns(m.Data),
	}
}

func (m Message) ToColumns(data map[string]any) []Column {
	var columns []Column
	for k, v := range data {
		columns = append(columns, Column{
			Name:  k,
			Value: v,
		})
	}
	return columns
}
