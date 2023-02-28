package aisera

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func ToJSONReader(val any) io.Reader {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(val)
	if err != nil {
		logger.Error("error encoding struct", err)
	}
	return &buf
}

func JSONReaderToVal(r io.Reader, val any) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	err = json.Unmarshal(data, val)
	if err != nil {
		return fmt.Errorf("error decoding the response: %q: %w", string(data), err)
	}
	return nil
}
