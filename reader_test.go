package orderedmap

import (
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	m := New()
	m.Set("name", "test")
	m.Set("value", 42.0)

	data, err := io.ReadAll(m.Reader())
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	want := `{"name":"test","value":42}`
	if got := string(data); got != want {
		t.Fatalf("Reader() = %s, want %s", got, want)
	}
}

func TestReaderEmpty(t *testing.T) {
	m := New()
	data, err := io.ReadAll(m.Reader())
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	if got := string(data); got != "{}" {
		t.Fatalf("Reader() = %s, want {}", got)
	}
}
