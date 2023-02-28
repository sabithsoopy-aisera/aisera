package aisera

import (
	"io"
)

type SortCriteria struct {
	Field     string `json:"field"`
	Ascending bool   `json:"ascending"`
}

type Filter struct {
	Fields          []string       `json:"fields"`
	SortCriteria    []SortCriteria `json:"sortCriteria"`
	MaxCount        int            `json:"maxCount"`
	NeedsTotalCount bool           `json:"needsTotalCount"`
	Offset          int            `json:"offset"`
}

func (f Filter) JSONReader() io.Reader {
	return ToJSONReader(f)
}
