package staticfile

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/files/stateloader"

	"github.com/getsyncer/syncer-core/drift/syncers/templatefiles/templatemutator"

	"github.com/getsyncer/syncer-core/drift"

	"github.com/getsyncer/syncer-core/config"

	"github.com/getsyncer/syncer-core/files"
)

type Config struct {
	// Filename is the name of the file to create
	Filename string `yaml:"filename"`
	// Content is the content of the file to create
	Content string `yaml:"content"`
}

type CreateFilesystem interface {
	Changes(ctx context.Context) (files.System[*files.StateWithChangeReason], error)
}

func (c Config) Changes(_ context.Context) (files.System[*files.StateWithChangeReason], error) {
	var ret files.System[*files.StateWithChangeReason]
	if c.Filename == "" {
		return ret, fmt.Errorf("filename is required")
	}
	path := c.Filename
	if err := ret.Add(files.Path(c.Filename), &files.StateWithChangeReason{
		ChangeReason: &files.ChangeReason{
			Reason: "staticfile",
		},
		State: files.State{
			Mode:          0644,
			Contents:      []byte(c.Content),
			FileExistence: files.FileExistencePresent,
		},
	}); err != nil {
		return ret, fmt.Errorf("failed to add file %q: %w", path, err)
	}
	return ret, nil
}

// New returns a new staticfile syncer
func New(loader stateloader.StateLoader) *Syncer[Config] {
	return NewWithCustomLogic[Config]("staticfile", drift.PriorityNormal, loader)
}

func NewWithCustomLogic[T CreateFilesystem](name config.Name, priority drift.Priority, loader stateloader.StateLoader) *Syncer[T] {
	return &Syncer[T]{
		name:     name,
		priority: priority,
		loader:   loader,
	}
}

type Syncer[T CreateFilesystem] struct {
	name     config.Name
	priority drift.Priority
	mutators templatemutator.MutatorList[T]
	loader   stateloader.StateLoader
}

// DetectDrift returns a system containing the file specified in the config
func (f *Syncer[T]) DetectDrift(ctx context.Context, runData *drift.RunData) (*files.System[*files.StateWithChangeReason], error) {
	var cfg T
	if err := runData.RunConfig.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staticfile config: %w", err)
	}
	cfg, err := f.mutators.Mutate(ctx, runData, f.loader, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to mutate config: %w", err)
	}
	ret, err := cfg.Changes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to apply changes")
	}
	return &ret, nil
}

func (f *Syncer[T]) AddMutator(mutator templatemutator.Mutator[T]) {
	f.mutators.AddMutator(mutator)
}

func (f *Syncer[T]) Name() config.Name {
	return f.name
}

func (f *Syncer[T]) Priority() drift.Priority {
	return f.priority
}

var _ drift.Detector = &Syncer[Config]{}
