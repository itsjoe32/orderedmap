// Package orderedmap provides an ordered map that preserves insertion order.
//
// Keys are strings and values are any type. All operations (Get, Set, Delete)
// run in amortized O(1) time. Iteration via [OrderedMap.Range] visits entries
// in insertion order and is safe to call even if the callback deletes entries.
//
// The zero value is ready to use; calling [New] is optional.
//
// OrderedMap implements [encoding/json.Marshaler], [encoding/json.Unmarshaler],
// [gopkg.in/yaml.v3.Marshaler], and [gopkg.in/yaml.v3.Unmarshaler].
// Nested JSON/YAML objects are decoded as *OrderedMap to preserve ordering
// recursively.
package orderedmap

// OrderedMap is a map with string keys that maintains insertion order.
// It is backed by a hash map for O(1) lookups and a doubly linked list
// embedded in a slice for ordered traversal. Deleted slots are recycled
// through a free list.
//
// The zero value is an empty map ready to use.
type OrderedMap struct {
	index   map[string]uint32
	entries []entry
	head    uint32
	tail    uint32
	free    uint32
}

type entry struct {
	key      string
	value    any
	prev     uint32
	next     uint32
	freeNext uint32
	used     bool
}

// New returns an initialized, empty OrderedMap.
func New() *OrderedMap {
	m := &OrderedMap{}
	m.init()
	return m
}

func (m *OrderedMap) init() {
	if m == nil || m.index != nil {
		return
	}
	m.index = make(map[string]uint32)
	m.entries = make([]entry, 1)
}

// Get returns the value associated with key and true, or nil and false
// if the key is not present. It is safe to call on a nil or zero-value map.
func (m *OrderedMap) Get(key string) (any, bool) {
	if m == nil || m.index == nil {
		return nil, false
	}
	i, ok := m.index[key]
	if !ok {
		return nil, false
	}
	return m.entries[i].value, true
}

// Set adds or updates the entry for key. It reports whether the key was
// newly inserted (true) or an existing value was updated (false).
// A new key is appended to the end of the iteration order.
// It is safe to call on a zero-value map. Calling it on a nil pointer is a no-op.
func (m *OrderedMap) Set(key string, value any) bool {
	if m == nil {
		return false
	}
	m.init()

	if i, ok := m.index[key]; ok {
		m.entries[i].value = value
		return false
	}

	var i uint32
	if m.free != 0 {
		i = m.free
		m.free = m.entries[i].freeNext
	} else {
		i = uint32(len(m.entries))
		m.entries = append(m.entries, entry{})
	}

	m.entries[i] = entry{
		key:   key,
		value: value,
		prev:  m.tail,
		used:  true,
	}
	if m.tail != 0 {
		m.entries[m.tail].next = i
	} else {
		m.head = i
	}
	m.tail = i
	m.index[key] = i
	return true
}

// Delete removes the entry for key. It reports whether the key was present.
// The deleted slot is added to an internal free list for reuse by future
// Set calls. It is safe to call on a nil or zero-value map.
func (m *OrderedMap) Delete(key string) bool {
	if m == nil || m.index == nil {
		return false
	}

	i, ok := m.index[key]
	if !ok {
		return false
	}

	e := &m.entries[i]
	if e.prev != 0 {
		m.entries[e.prev].next = e.next
	} else {
		m.head = e.next
	}
	if e.next != 0 {
		m.entries[e.next].prev = e.prev
	} else {
		m.tail = e.prev
	}

	delete(m.index, key)

	e.key = ""
	e.value = nil
	e.used = false
	e.freeNext = m.free
	m.free = i
	return true
}

// Len returns the number of entries in the map.
func (m *OrderedMap) Len() int {
	if m == nil || m.index == nil {
		return 0
	}
	return len(m.index)
}

// Range calls yield for each entry in insertion order. If yield returns
// false, iteration stops. It is safe to delete entries during iteration.
// Range is compatible with range-over-func (Go 1.23+).
func (m *OrderedMap) Range(yield func(string, any) bool) {
	if m == nil {
		return
	}
	for i := m.head; i != 0; {
		e := m.entries[i]
		next := e.next
		if e.used && !yield(e.key, e.value) {
			return
		}
		i = next
	}
}

// Keys returns all keys in insertion order.
func (m *OrderedMap) Keys() []string {
	keys := make([]string, 0, m.Len())
	m.Range(func(key string, _ any) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// Values returns all values in insertion order.
func (m *OrderedMap) Values() []any {
	values := make([]any, 0, m.Len())
	m.Range(func(_ string, value any) bool {
		values = append(values, value)
		return true
	})
	return values
}
