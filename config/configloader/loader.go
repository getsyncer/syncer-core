package configloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/getsyncer/syncer-core/drift"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/syncer/childrenregistry"
	"gopkg.in/yaml.v3"
)

type DefaultConfigLoader struct {
	childrenRegistry childrenregistry.ChildrenRegistry
}

func NewDefaultConfigLoader(childrenRegistry childrenregistry.ChildrenRegistry) *DefaultConfigLoader {
	return &DefaultConfigLoader{
		childrenRegistry: childrenRegistry,
	}
}

type ConfigLoader interface {
	LoadConfig(ctx context.Context, contents io.Reader) (*config.Root, error)
	FlattenChildren(ctx context.Context, root *config.Root) error
}

var _ ConfigLoader = &DefaultConfigLoader{}

func DefaultLocations() []string {
	return []string{
		filepath.Join(drift.DefaultSyncerDirectory, drift.DefaultSyncerConfigFileName),
		filepath.Join(drift.DefaultSyncerDirectory, "config.yaml"),
		drift.DefaultSyncerConfigFileName,
	}
}

func DefaultFindConfigFile(wd string) (string, error) {
	if wd == "" {
		var err error
		wd, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
	}
	possibleLocations := DefaultLocations()
	for _, loc := range possibleLocations {
		fileLoc := filepath.Join(wd, loc)
		if _, err := os.Stat(fileLoc); err == nil {
			return loc, nil
		}
	}
	return "", fmt.Errorf("no config file found inside %s", wd)
}

func ConfigFromFile(ctx context.Context, loader ConfigLoader) (*config.Root, error) {
	f, err := DefaultFindConfigFile("")
	if err != nil {
		return nil, fmt.Errorf("failed to find config file: %w", err)
	}
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return loader.LoadConfig(ctx, bytes.NewReader(b))
}

func (c *DefaultConfigLoader) LoadConfig(_ context.Context, contents io.Reader) (*config.Root, error) {
	var root config.Root
	dec := yaml.NewDecoder(contents)
	if err := dec.Decode(&root); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	if root.Version != 1 {
		return nil, fmt.Errorf("unknown config version: %d", root.Version)
	}
	return &root, nil
}

func (c *DefaultConfigLoader) FlattenChildren(_ context.Context, root *config.Root) error {
	// Now process children
	processedChildrenSources := make(map[string]struct{})
	i := 0
	for len(root.Children) > 0 {
		i++
		if i > 1000 {
			return fmt.Errorf("too many children: forever loop protection")
		}
		originalChildren := root.Children
		root.Children = nil
		for _, child := range originalChildren {
			content, exists := c.childrenRegistry.Get(config.Name(child.Source))
			if !exists {
				return fmt.Errorf("unknown child: %s", child.Source)
			}
			if _, exists := processedChildrenSources[child.Source]; exists {
				continue
			}
			processedChildrenSources[child.Source] = struct{}{}
			// Now unmarshal the child
			var childRoot config.Root
			if err := yaml.Unmarshal(content.Content, &childRoot); err != nil {
				return fmt.Errorf("failed to unmarshal child: %w", err)
			}
			if err := root.Merge(&childRoot); err != nil {
				return fmt.Errorf("failed to merge child: %w", err)
			}
		}
	}
	return nil
}
