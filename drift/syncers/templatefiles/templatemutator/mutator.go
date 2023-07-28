package templatemutator

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files/stateloader"
)

type Mutator[T config.TemplateConfig] interface {
	Mutate(ctx context.Context, runData *drift.RunData, loader stateloader.StateLoader, cfg T) (T, error)
}

type MutatorFunc[T config.TemplateConfig] func(ctx context.Context, runData *drift.RunData, loader stateloader.StateLoader, cfg T) (T, error)

func (c MutatorFunc[T]) Mutate(ctx context.Context, runData *drift.RunData, loader stateloader.StateLoader, cfg T) (T, error) {
	return c(ctx, runData, loader, cfg)
}

func SimpleMutator[T config.TemplateConfig](updateFunc func(T) (T, error)) Mutator[T] {
	return MutatorFunc[T](func(_ context.Context, _ *drift.RunData, _ stateloader.StateLoader, cfg T) (T, error) {
		return updateFunc(cfg)
	})
}

type Mutatable[T config.TemplateConfig] interface {
	AddMutator(mutator Mutator[T])
}

type MutatorList[T config.TemplateConfig] struct {
	mutators []Mutator[T]
}

var _ Mutator[config.TemplateConfig] = &MutatorList[config.TemplateConfig]{}
var _ Mutatable[config.TemplateConfig] = &MutatorList[config.TemplateConfig]{}

func (m *MutatorList[T]) AddMutator(mutator Mutator[T]) {
	m.mutators = append(m.mutators, mutator)
}

func (m *MutatorList[T]) Mutate(ctx context.Context, runData *drift.RunData, loader stateloader.StateLoader, cfg T) (T, error) {
	for _, mutator := range m.mutators {
		var err error
		cfg, err = mutator.Mutate(ctx, runData, loader, cfg)
		if err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func addMutator[T config.TemplateConfig](r drift.Registry, name config.Name, mutator Mutator[T]) error {
	s, ok := r.Get(name)
	if !ok {
		return fmt.Errorf("syncer %s not found", name)
	}
	asMutatable, ok := s.(Mutatable[T])
	if !ok {
		return fmt.Errorf("syncer %s is not mutatable", name)
	}
	asMutatable.AddMutator(mutator)
	return nil
}

type SetupMutator[T config.TemplateConfig] struct {
	Mutator Mutator[T]
	Name    config.Name
}

func (s *SetupMutator[T]) Setup(_ context.Context, runData *drift.RunData) error {
	if err := addMutator[T](runData.Registry, s.Name, s.Mutator); err != nil {
		return fmt.Errorf("unable to add mutator: %w", err)
	}
	return nil
}
