package templatemutator

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/stateloader"
)

type PostGenConfigMutator[T config.TemplateConfig] struct {
	ToMutate           config.Name
	TemplateName       string
	PostGenMutatorFunc PostGenMutatorFunc[T]
}

type PostGenMutatorFunc[T config.TemplateConfig] func(ctx context.Context, renderedTemplate string, cfg T) (T, error)

func (p *PostGenConfigMutator[T]) makeMutator(renderedTemplate string) Mutator[T] {
	return MutatorFunc[T](func(ctx context.Context, runData *drift.RunData, loader stateloader.StateLoader, cfg T) (T, error) {
		return p.PostGenMutatorFunc(ctx, renderedTemplate, cfg)
	})
}

func (p *PostGenConfigMutator[T]) PostGenProcess(_ context.Context, fs *files.System[*files.StateWithChangeReason], runData *drift.RunData) error {
	s, exists := fs.Remove(files.Path(p.TemplateName))
	if !exists {
		return fmt.Errorf("unable to find template file %s", p.TemplateName)
	}
	if err := addMutator(runData.Registry, p.ToMutate, p.makeMutator(string(s.State.Contents))); err != nil {
		return fmt.Errorf("unable to add mutator: %w", err)
	}
	return nil
}
