package config

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Root struct {
	Children []Logic `yaml:"children,omitempty"`
	Version  int     `yaml:"version"`
	Config   Dynamic `yaml:"config,omitempty"`
	Logic    []Logic `yaml:"logic,omitempty"`
	Syncs    []Syncs `yaml:"syncs,omitempty"`
}

func (r *Root) AsYaml() (string, error) {
	var into bytes.Buffer
	enc := yaml.NewEncoder(&into)
	enc.SetIndent(2)
	if err := enc.Encode(r); err != nil {
		return "", fmt.Errorf("error encoding yaml: %w", err)
	}
	return into.String(), nil
}

func (r *Root) Merge(other *Root) error {
	if other == nil {
		return nil
	}
	if other.Version != r.Version {
		return fmt.Errorf("cannot merge different config versions: %d vs %d", r.Version, other.Version)
	}
	if other.Children != nil {
		r.Children = append(r.Children, other.Children...)
		r.Children = removeDuplicate(r.Children)
	}
	if other.Logic != nil {
		r.Logic = append(r.Logic, other.Logic...)
		r.Logic = removeDuplicate(r.Logic)
	}
	if other.Syncs != nil {
		r.Syncs = append(r.Syncs, other.Syncs...)
	}
	if err := r.Config.Merge(other.Config); err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}
	return nil
}

func removeDuplicate[T comparable](xs []T) []T {
	var out []T
	seen := make(map[T]struct{})
	for _, x := range xs {
		if _, ok := seen[x]; !ok {
			seen[x] = struct{}{}
			out = append(out, x)
		}
	}
	return out
}

type TemplateConfig interface{}
