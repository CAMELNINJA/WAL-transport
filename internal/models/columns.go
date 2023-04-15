package models

import "fmt"

type ColumnInfo struct {
	TableName  string
	ColumnName string
	DataType   string
	ColumnNum  int
}

type Columns []*ColumnInfo

func (c *ColumnInfo) String() string {
	return fmt.Sprintf("TableName %s ColumnName %s DataType %s", c.TableName, c.ColumnName, c.DataType)
}
