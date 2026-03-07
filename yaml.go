package orderedmap

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// NewFromYAML parses a YAML mapping and returns an OrderedMap with keys in
// the order they appear in the input. Nested mappings are decoded as *OrderedMap.
func NewFromYAML(data []byte) (*OrderedMap, error) {
	m := New()
	if err := yaml.Unmarshal(data, m); err != nil {
		return nil, err
	}
	return m, nil
}

// MarshalYAML implements yaml.Marshaler.
func (m *OrderedMap) MarshalYAML() (any, error) {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}
	var marshalErr error
	m.Range(func(key string, value any) bool {
		keyNode := &yaml.Node{}
		if err := keyNode.Encode(key); err != nil {
			marshalErr = err
			return false
		}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(value); err != nil {
			marshalErr = err
			return false
		}
		node.Content = append(node.Content, keyNode, valueNode)
		return true
	})
	if marshalErr != nil {
		return nil, marshalErr
	}
	return node, nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (m *OrderedMap) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("orderedmap: expected mapping node, got kind %d", value.Kind)
	}
	*m = OrderedMap{}
	m.init()
	for i := 0; i < len(value.Content)-1; i += 2 {
		keyNode := value.Content[i]
		valNode := value.Content[i+1]
		var key string
		if err := keyNode.Decode(&key); err != nil {
			return err
		}
		val, err := decodeYAMLNode(valNode)
		if err != nil {
			return err
		}
		m.Set(key, val)
	}
	return nil
}

func decodeYAMLNode(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.MappingNode:
		m := New()
		if err := node.Decode(m); err != nil {
			return nil, err
		}
		return m, nil
	case yaml.SequenceNode:
		result := make([]any, len(node.Content))
		for i, child := range node.Content {
			v, err := decodeYAMLNode(child)
			if err != nil {
				return nil, err
			}
			result[i] = v
		}
		return result, nil
	default:
		var v any
		if err := node.Decode(&v); err != nil {
			return nil, err
		}
		return v, nil
	}
}
