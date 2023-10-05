package config

import (
	"bytes"
	"errors"
	"fmt"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Dynamic struct {
	root *yaml.Node
}

func (r Dynamic) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	if r.root == nil {
		return nil
	}
	for i := 0; i < len(r.root.Content); i += 2 {
		if i+1 >= len(r.root.Content) {
			return errors.New("odd number of nodes in yaml")
		}
		key := r.root.Content[i]
		value := r.root.Content[i+1]
		if key.Kind == yaml.ScalarNode {
			encoder.AddString(key.Value, value.Value)
		}
	}
	return nil
}

var _ zapcore.ObjectMarshaler = &Dynamic{}
var _ yaml.Marshaler = &Dynamic{}
var _ yaml.Unmarshaler = &Dynamic{}

func (r *Dynamic) Decode(into interface{}) error {
	if r.root == nil {
		return nil
	}
	return r.root.Decode(into)
}

func (r Dynamic) MarshalYAML() (interface{}, error) {
	// Note: Because of https://github.com/golang/go/issues/22967 we can't use a pointer receiver here.
	return r.root, nil
}

func (r *Dynamic) SetValueIfEmpty(key string, value string) {
	// if r.root is nil, set it to a new yaml.Node
	if r.root == nil {
		r.root = &yaml.Node{
			Kind: yaml.MappingNode,
		}
	}

	// if r.root[key] is empty, set it to value
	for i := 0; i < len(r.root.Content); i += 2 {
		if i+1 >= len(r.root.Content) {
			return
		}
		if r.root.Content[i].Value == key {
			if r.root.Content[i+1].Value == "" {
				r.root.Content[i+1].Value = value
			}
			return
		}
	}
	// if r.root[key] is not empty, append key and value to r.root
	r.root.Content = append(r.root.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: key,
	}, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: value,
	})
}

func (r *Dynamic) UnmarshalYAML(value *yaml.Node) error {
	r.root = value
	return nil
}

func (r *Dynamic) Merge(other Dynamic) error {
	if other.root == nil {
		return nil
	}
	if r.root == nil {
		varCopy := *other.root
		r.root = &varCopy
		return nil
	}
	if err := recursiveMerge(other.root, r.root); err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}
	return nil
}

func (r *Dynamic) AsYaml() (string, error) {
	var into bytes.Buffer
	enc := yaml.NewEncoder(&into)
	enc.SetIndent(2)
	if err := enc.Encode(r); err != nil {
		return "", fmt.Errorf("error encoding yaml: %w", err)
	}
	return into.String(), nil
}

func nodesEqual(l, r *yaml.Node) bool {
	if l.Kind == yaml.ScalarNode && r.Kind == yaml.ScalarNode {
		return l.Value == r.Value
	}
	panic("equals on non-scalars not implemented!")
}

// https://stackoverflow.com/questions/65768861/read-and-merge-two-yaml-files-dynamically-and-or-recursively
func recursiveMerge(from, into *yaml.Node) error {
	if from.Kind != into.Kind {
		return fmt.Errorf("cannot merge nodes of different kinds: from=%d vs into=%d", from.Kind, into.Kind)
	}
	switch from.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(from.Content); i += 2 {
			found := false
			for j := 0; j < len(into.Content); j += 2 {
				if nodesEqual(from.Content[i], into.Content[j]) {
					found = true
					if err := recursiveMerge(from.Content[i+1], into.Content[j+1]); err != nil {
						return fmt.Errorf("failed to merge mapping node at %s: %w", from.Content[i].Value, err)
					}
					break
				}
			}
			if !found {
				into.Content = append(into.Content, from.Content[i:i+2]...)
			}
		}
	case yaml.SequenceNode:
		into.Content = append(into.Content, from.Content...)
	case yaml.DocumentNode:
		if err := recursiveMerge(from.Content[0], into.Content[0]); err != nil {
			return fmt.Errorf("failed to merge document node: %w", err)
		}
	case yaml.ScalarNode:
		// Skip scalars (It's like merging one value having "v3" into another value that already has "v4".  Just keep the "v4" value.)
	default:
		return fmt.Errorf("can only merge mapping and sequence nodes, not %d", from.Kind)
	}
	return nil
}
