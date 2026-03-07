package orderedmap

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		give func() *OrderedMap
		want string
	}{
		{
			name: "EmptyMap",
			give: func() *OrderedMap { return New() },
			want: `{}`,
		},
		{
			name: "PreservesOrder",
			give: func() *OrderedMap {
				m := New()
				m.Set("z", 1.0)
				m.Set("a", "hello")
				m.Set("m", true)
				return m
			},
			want: `{"z":1,"a":"hello","m":true}`,
		},
		{
			name: "NestedOrderedMap",
			give: func() *OrderedMap {
				inner := New()
				inner.Set("y", 2.0)
				inner.Set("x", 1.0)
				m := New()
				m.Set("nested", inner)
				m.Set("top", "val")
				return m
			},
			want: `{"nested":{"y":2,"x":1},"top":"val"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.give())
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		give     string
		wantKeys []string
	}{
		{"PreservesOrder", `{"z":1,"a":"hello","m":true}`, []string{"z", "a", "m"}},
		{"NestedObject", `{"outer":{"b":2,"a":1}}`, []string{"outer"}},
		{"WithArray", `{"items":[1,2,3]}`, []string{"items"}},
		{"WithNull", `{"a":null}`, []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if err := json.Unmarshal([]byte(tt.give), m); err != nil {
				t.Fatalf("Unmarshal() error: %v", err)
			}
			if got := m.Keys(); !slices.Equal(got, tt.wantKeys) {
				t.Fatalf("Keys() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}

func TestUnmarshalJSONNestedObjectsAreOrderedMaps(t *testing.T) {
	m, err := NewFromJSON([]byte(`{"outer":{"b":2,"a":1}}`))
	if err != nil {
		t.Fatal(err)
	}
	v, _ := m.Get("outer")
	nested, ok := v.(*OrderedMap)
	if !ok {
		t.Fatalf("nested value type = %T, want *OrderedMap", v)
	}
	want := []string{"b", "a"}
	if got := nested.Keys(); !slices.Equal(got, want) {
		t.Fatalf("nested Keys() = %v, want %v", got, want)
	}
}

func TestUnmarshalJSONArraysWithNestedObjects(t *testing.T) {
	m, err := NewFromJSON([]byte(`{"list":[{"b":2,"a":1}]}`))
	if err != nil {
		t.Fatal(err)
	}
	v, _ := m.Get("list")
	arr := v.([]any)
	elem := arr[0].(*OrderedMap)
	want := []string{"b", "a"}
	if got := elem.Keys(); !slices.Equal(got, want) {
		t.Fatalf("elem Keys() = %v, want %v", got, want)
	}
}

func TestUnmarshalJSONInvalidInput(t *testing.T) {
	m := New()
	if err := m.UnmarshalJSON([]byte(`[]`)); err == nil {
		t.Fatal("UnmarshalJSON() error = nil, want error")
	}
}

func TestUnmarshalJSONRejectsTrailingData(t *testing.T) {
	m := New()
	if err := m.UnmarshalJSON([]byte(`{"a":1}{"b":2}`)); err == nil {
		t.Fatal("UnmarshalJSON() error = nil, want error")
	}
}

func TestNewFromJSON(t *testing.T) {
	m, err := NewFromJSON([]byte(`{"x":1,"y":2}`))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"x", "y"}
	if got := m.Keys(); !slices.Equal(got, want) {
		t.Fatalf("Keys() = %v, want %v", got, want)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	m := New()
	m.Set("c", 3.0)
	m.Set("a", 1.0)
	m.Set("b", 2.0)

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	m2 := New()
	if err := json.Unmarshal(data, m2); err != nil {
		t.Fatal(err)
	}
	if got := m2.Keys(); !slices.Equal(got, m.Keys()) {
		t.Fatalf("Keys() = %v, want %v", got, m.Keys())
	}
}
