package orderedmap

import (
	"bytes"
	"io"
)

// Reader returns an io.Reader that yields the JSON encoding of the map.
// This is useful for passing the map directly as an HTTP request body:
//
//	req, err := http.NewRequest("POST", url, m.Reader())
func (m *OrderedMap) Reader() io.Reader {
	data, err := m.MarshalJSON()
	if err != nil {
		return &errReader{err: err}
	}
	return bytes.NewReader(data)
}

type errReader struct {
	err error
}

func (r *errReader) Read([]byte) (int, error) {
	return 0, r.err
}
