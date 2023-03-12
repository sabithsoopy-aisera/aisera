package aisera

import (
	"io"
)

type SortCriteria struct {
	Field     string `json:"field,omitempty"`
	Ascending bool   `json:"ascending,omitempty"`
}

type Criteria struct {
	Field    []string `json:"field,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Operand  string   `json:"operand,omitempty"`
	CanImply bool     `json:"canImply,omitempty"`
}

type Filter struct {
	EntityId        int            `json:"entityId,omitempty"`
	Fields          []string       `json:"fields,omitempty"`
	Criteria        []Criteria     `json:"criteria,omitempty"`
	SortCriteria    []SortCriteria `json:"sortCriteria,omitempty"`
	MaxCount        int            `json:"maxCount,omitempty"`
	NeedsTotalCount bool           `json:"needsTotalCount,omitempty"`
	Offset          int            `json:"offset,omitempty"`
}

func (f Filter) JSONReader() io.Reader {
	return ToJSONReader(f)
}
