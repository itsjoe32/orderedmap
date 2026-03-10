# orderedmap

A lightweight ordered map for Go that preserves string-key insertion order, keeps Get/Set/Delete O(1), and retains key order when marshaling JSON or YAML.

Documentation: https://pkg.go.dev/github.com/itsjoe32/orderedmap

## Install

```
go get github.com/itsjoe32/orderedmap
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/itsjoe32/orderedmap"
)

func main() {
	m := orderedmap.New()
	m.Set("z", 1)
	m.Set("a", 2)
	m.Set("m", 3)

	fmt.Println(m.Keys())   // [z a m]
	fmt.Println(m.Values()) // [1 2 3]

	value, ok := m.Get("a")
	fmt.Println(value, ok) // 2 true

	m.Delete("a")
	fmt.Println(m.Len()) // 2

	for k, v := range m.Range {
		fmt.Println(k, v)
	}
}
```

The zero value is ready to use:

```go
var m orderedmap.OrderedMap
m.Set("key", "value")
```

## JSON

Implements `json.Marshaler` and `json.Unmarshaler`. Key order is preserved in both directions. Nested objects are decoded as `*OrderedMap`.

```go
jsonMap := orderedmap.New()
jsonMap.Set("z", 1)
jsonMap.Set("a", 2)

data, _ := json.Marshal(jsonMap)
// {"z":1,"a":2}

m, err := orderedmap.NewFromJSON([]byte(`{"z":1,"a":2}`))
```

## YAML

Implements `yaml.Marshaler` and `yaml.Unmarshaler` ([gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)). Nested mappings are decoded as `*OrderedMap`.

```go
yamlMap := orderedmap.New()
yamlMap.Set("z", 1)
yamlMap.Set("a", 2)

data, _ := yaml.Marshal(yamlMap)
// z: 1
// a: 2

m, err := orderedmap.NewFromYAML([]byte("z: 1\na: 2\n"))
```

## HTTP Request Body

`Reader()` returns an `io.Reader` with the JSON-encoded map, ready to use as an HTTP request body:

```go
m := orderedmap.New()
m.Set("name", "test")
m.Set("value", 42)

req, err := http.NewRequest("POST", "https://api.example.com/data", m.Reader())
req.Header.Set("Content-Type", "application/json")
```

## Constraints

- Keys are strings.
- Iteration preserves insertion order.
- Updating an existing key changes its value but does not move its position.
- Nested JSON objects and YAML mappings are decoded as `*OrderedMap`.

## Compatibility

- The `for k, v := range m.Range` syntax requires Go 1.23 or newer.
- On older Go versions, iterate with the callback form:

```go
m.Range(func(k string, v any) bool {
	fmt.Println(k, v)
	return true
})
```

## API

| Method | Description |
|---|---|
| `New()` | Returns an initialized, empty map |
| `Get(key)` | Returns value and existence bool |
| `Set(key, value)` | Inserts or updates; returns true if key is new |
| `Delete(key)` | Removes key; returns true if key existed |
| `Len()` | Returns number of entries |
| `Keys()` | Returns keys in insertion order |
| `Values()` | Returns values in insertion order |
| `Range(yield)` | Iterates in order; safe to delete during iteration |
| `NewFromJSON(data)` | Parses JSON into an OrderedMap |
| `NewFromYAML(data)` | Parses YAML into an OrderedMap |
| `Reader()` | Returns an `io.Reader` with the JSON-encoded map |
