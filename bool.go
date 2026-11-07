package plex

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Bool is a boolean that also decodes from the numbers 0 and 1.
//
// Plex Media Server converts its XML attributes to JSON and does not always
// emit the type its specification declares: some fields documented as
// boolean are sent as 0 or 1. Fields known to do this use Bool instead of
// bool. Bool encodes as a JSON boolean.
type Bool bool

// UnmarshalJSON implements [json.Unmarshaler]. It accepts true, false, 0, 1,
// and the same values as JSON strings. A JSON null leaves b unchanged.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		return nil
	}
	if s, err := unquote(data); err == nil {
		data = s
	}
	switch string(data) {
	case "true", "1":
		*b = true
	case "false", "0":
		*b = false
	default:
		return fmt.Errorf("plex: cannot decode %s into Bool", data)
	}
	return nil
}

// unquote returns the contents of a JSON string, or an error if data is not
// a JSON string.
func unquote(data []byte) ([]byte, error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return []byte(s), nil
}
