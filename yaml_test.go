package orderedmap

import (
	"slices"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		give func() *OrderedMap
		want string
	}{
		{
			name: "EmptyMap",
			give: func() *OrderedMap { return New() },
			want: "{}\n",
		},
		{
			name: "PreservesOrder",
			give: func() *OrderedMap {
				m := New()
				m.Set("z", 1)
				m.Set("a", "hello")
				m.Set("m", true)
				return m
			},
			want: "z: 1\na: hello\nm: true\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := yaml.Marshal(tt.give())
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal() =\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		give     string
		wantKeys []string
	}{
		{"PreservesOrder", "z: 1\na: hello\nm: true\n", []string{"z", "a", "m"}},
		{"NestedMapping", "outer:\n  b: 2\n  a: 1\n", []string{"outer"}},
		{"WithSequence", "items:\n  - 1\n  - 2\n  - 3\n", []string{"items"}},
		{"WithNull", "a: null\n", []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if err := yaml.Unmarshal([]byte(tt.give), m); err != nil {
				t.Fatalf("Unmarshal() error: %v", err)
			}
			if got := m.Keys(); !slices.Equal(got, tt.wantKeys) {
				t.Fatalf("Keys() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}

func TestUnmarshalYAMLNestedMappingsAreOrderedMaps(t *testing.T) {
	m, err := NewFromYAML([]byte("outer:\n  b: 2\n  a: 1\n"))
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

func TestUnmarshalYAMLSequencesWithNestedMappings(t *testing.T) {
	m, err := NewFromYAML([]byte("list:\n  - b: 2\n    a: 1\n"))
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

func TestUnmarshalYAMLInvalidInput(t *testing.T) {
	m := New()
	if err := yaml.Unmarshal([]byte("- item1\n- item2\n"), m); err == nil {
		t.Fatal("Unmarshal() error = nil, want error")
	}
}

func TestNewFromYAML(t *testing.T) {
	m, err := NewFromYAML([]byte("x: 1\ny: 2\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"x", "y"}
	if got := m.Keys(); !slices.Equal(got, want) {
		t.Fatalf("Keys() = %v, want %v", got, want)
	}
}

func TestYAMLRoundTrip(t *testing.T) {
	m := New()
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)

	data, err := yaml.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	m2, err := NewFromYAML(data)
	if err != nil {
		t.Fatal(err)
	}
	if got := m2.Keys(); !slices.Equal(got, m.Keys()) {
		t.Fatalf("Keys() = %v, want %v", got, m.Keys())
	}
}
