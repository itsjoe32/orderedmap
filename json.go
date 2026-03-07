package orderedmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// NewFromJSON parses a JSON object and returns an OrderedMap with keys in
// the order they appear in the input. Nested objects are decoded as *OrderedMap.
func NewFromJSON(data []byte) (*OrderedMap, error) {
	m := New()
	if err := m.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return m, nil
}

// MarshalJSON implements json.Marshaler.
func (m *OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	var marshalErr error
	m.Range(func(key string, value any) bool {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		k, err := json.Marshal(key)
		if err != nil {
			marshalErr = err
			return false
		}
		buf.Write(k)
		buf.WriteByte(':')
		v, err := json.Marshal(value)
		if err != nil {
			marshalErr = err
			return false
		}
		buf.Write(v)
		return true
	})
	if marshalErr != nil {
		return nil, marshalErr
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *OrderedMap) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := t.(json.Delim)
	if !ok || delim != '{' {
		return fmt.Errorf("orderedmap: expected '{', got %v", t)
	}
	*m = OrderedMap{}
	m.init()
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return err
		}
		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("orderedmap: expected string key, got %v", t)
		}
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return err
		}
		value, err := decodeJSONValue(raw)
		if err != nil {
			return err
		}
		m.Set(key, value)
	}
	t, err = dec.Token()
	if err != nil {
		return err
	}
	delim, ok = t.(json.Delim)
	if !ok || delim != '}' {
		return fmt.Errorf("orderedmap: expected '}', got %v", t)
	}
	if _, err := dec.Token(); err != io.EOF {
		if err == nil {
			return fmt.Errorf("orderedmap: unexpected trailing data after object")
		}
		return err
	}
	return nil
}

func decodeJSONValue(raw json.RawMessage) (any, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("orderedmap: empty JSON value")
	}
	switch trimmed[0] {
	case '{':
		m := New()
		if err := m.UnmarshalJSON(trimmed); err != nil {
			return nil, err
		}
		return m, nil
	case '[':
		var arr []json.RawMessage
		if err := json.Unmarshal(trimmed, &arr); err != nil {
			return nil, err
		}
		result := make([]any, len(arr))
		for i, item := range arr {
			v, err := decodeJSONValue(item)
			if err != nil {
				return nil, err
			}
			result[i] = v
		}
		return result, nil
	default:
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		return v, nil
	}
}
