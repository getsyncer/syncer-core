package templatefiles

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files/stateloader"
)

type GenericConfigMutator[T config.TemplateConfig] struct {
	TemplateName string
	TemplateStr  string
	MutateFunc   func(ctx context.Context, renderedTemplate string, cfg T) (T, error)
}

func (g *GenericConfigMutator[T]) Mutate(ctx context.Context, runData *drift.RunData, _ stateloader.StateLoader, cfg T) (T, error) {
	updatedBuildGoLib, err := newTemplate(g.TemplateName, g.TemplateStr)
	if err != nil {
		return cfg, fmt.Errorf("unable to parse template: %w", err)
	}
	res, err := executeTemplateOnConfig(ctx, runData, cfg, updatedBuildGoLib)
	if err != nil {
		return cfg, fmt.Errorf("unable to execute template: %w", err)
	}
	return g.MutateFunc(ctx, res, cfg)
}
