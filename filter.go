package aisera

import (
	"io"
)

type SortCriteria struct {
	Field     string `json:"field,omitempty"`
	Ascending bool   `json:"ascending,omitempty"`
}

type Filter struct {
	Fields          []string       `json:"fields,omitempty"`
	SortCriteria    []SortCriteria `json:"sortCriteria,omitempty"`
	MaxCount        int            `json:"maxCount,omitempty"`
	NeedsTotalCount bool           `json:"needsTotalCount,omitempty"`
	Offset          int            `json:"offset,omitempty"`
}

func (f Filter) JSONReader() io.Reader {
	return ToJSONReader(f)
}
