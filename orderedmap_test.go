package orderedmap

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", m.Len())
	}
}

func TestGet(t *testing.T) {
	m := New()
	m.Set("foo", "bar")

	tests := []struct {
		name   string
		give   string
		wantV  any
		wantOK bool
	}{
		{"ReturnsValueForKey", "foo", "bar", true},
		{"ReturnsNotOKIfKeyDoesntExist", "missing", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := m.Get(tt.give)
			if ok != tt.wantOK || v != tt.wantV {
				t.Fatalf("Get(%q) = (%v, %v), want (%v, %v)", tt.give, v, ok, tt.wantV, tt.wantOK)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name    string
		give    []string
		wantNew bool
	}{
		{"ReturnsTrueIfKeyIsNew", []string{"foo"}, true},
		{"ReturnsFalseIfKeyExists", []string{"foo", "foo"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			var got bool
			for _, k := range tt.give {
				got = m.Set(k, k)
			}
			if got != tt.wantNew {
				t.Fatalf("Set() = %v, want %v", got, tt.wantNew)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		give     []string
		delete   string
		wantOK   bool
		wantKeys []string
	}{
		{"ReturnsFalseIfKeyDoesntExist", []string{}, "foo", false, []string{}},
		{"ReturnsTrueIfKeyExists", []string{"foo"}, "foo", true, []string{}},
		{"DeleteIsIsolated", []string{"foo", "bar"}, "foo", true, []string{"bar"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			for _, k := range tt.give {
				m.Set(k, k)
			}
			if got := m.Delete(tt.delete); got != tt.wantOK {
				t.Fatalf("Delete(%q) = %v, want %v", tt.delete, got, tt.wantOK)
			}
			if got := m.Keys(); !slices.Equal(got, tt.wantKeys) {
				t.Fatalf("Keys() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name string
		give []string
		want int
	}{
		{"EmptyMapIsZero", []string{}, 0},
		{"SingleElement", []string{"a"}, 1},
		{"ThreeElements", []string{"a", "b", "c"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			for _, k := range tt.give {
				m.Set(k, k)
			}
			if got := m.Len(); got != tt.want {
				t.Fatalf("Len() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		give []string
		want []string
	}{
		{"EmptyMap", []string{}, []string{}},
		{"RetainsOrder", []string{"z", "a", "m"}, []string{"z", "a", "m"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			for _, k := range tt.give {
				m.Set(k, k)
			}
			got := m.Keys()
			if len(got) != len(tt.want) || (len(got) > 0 && !slices.Equal(got, tt.want)) {
				t.Fatalf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	m := New()
	m.Set("a", 1)
	m.Set("b", "two")
	m.Set("c", 3.0)

	want := []any{1, "two", 3.0}
	got := m.Values()
	if len(got) != len(want) {
		t.Fatalf("len(Values()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Values()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestRangeIteratesInOrder(t *testing.T) {
	m := New()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	var keys []string
	m.Range(func(k string, _ any) bool {
		keys = append(keys, k)
		return true
	})
	want := []string{"a", "b", "c"}
	if !slices.Equal(keys, want) {
		t.Fatalf("visited = %v, want %v", keys, want)
	}
}

func TestRangeEarlyStop(t *testing.T) {
	m := New()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	var keys []string
	m.Range(func(k string, _ any) bool {
		keys = append(keys, k)
		return k != "b"
	})
	want := []string{"a", "b"}
	if !slices.Equal(keys, want) {
		t.Fatalf("visited = %v, want %v", keys, want)
	}
}

func TestRangeSafeToDeleteDuringIteration(t *testing.T) {
	m := New()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Range(func(k string, _ any) bool {
		if k == "b" {
			m.Delete("b")
		}
		return true
	})
	if m.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", m.Len())
	}
}
