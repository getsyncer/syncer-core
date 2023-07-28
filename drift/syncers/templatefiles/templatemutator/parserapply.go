package templatemutator

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/stateloader"
)

type ParserApply[T config.TemplateConfig] interface {
	ApplyFileParse(ctx context.Context, res *stateloader.ParseResult, path files.Path, runData *drift.RunData, cfg T) (T, error)
}

type ParserApplyFunc[T config.TemplateConfig] func(ctx context.Context, res *stateloader.ParseResult, path files.Path, runData *drift.RunData, cfg T) (T, error)

func (p ParserApplyFunc[T]) ApplyFileParse(ctx context.Context, res *stateloader.ParseResult, path files.Path, runData *drift.RunData, cfg T) (T, error) {
	return p(ctx, res, path, runData, cfg)
}

func DefaultParseMutator[T ApplyableConfig[T]](path files.Path) *ParserMutator[T] {
	return &ParserMutator[T]{
		Path:  path,
		Conf:  stateloader.RecommendedNewlineSeparatedConfig(),
		Apply: ConfigApply[T](),
	}
}

type ParserMutator[T config.TemplateConfig] struct {
	Path  files.Path
	Conf  stateloader.ParseConfig
	Apply ParserApply[T]
}

var _ Mutator[config.TemplateConfig] = &ParserMutator[config.TemplateConfig]{}

func (p *ParserMutator[T]) Mutate(ctx context.Context, runData *drift.RunData, loader stateloader.StateLoader, cfg T) (T, error) {
	pr, err := stateloader.Parse(ctx, loader, p.Path, p.Conf)
	if err != nil {
		return cfg, fmt.Errorf("unable to parse: %w", err)
	}
	ret, err := p.Apply.ApplyFileParse(ctx, pr, p.Path, runData, cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to apply file parse: %w", err)
	}
	return ret, nil
}

type ApplyableConfig[T config.TemplateConfig] interface {
	ApplyParse(parse *stateloader.ParseResult) (T, error)
}

func ConfigApply[T ApplyableConfig[T]]() ParserApply[T] {
	return ParserApplyFunc[T](func(_ context.Context, res *stateloader.ParseResult, _ files.Path, _ *drift.RunData, cfg T) (T, error) {
		ret, err := cfg.ApplyParse(res)
		if err != nil {
			return cfg, fmt.Errorf("unable to apply parse: %w", err)
		}
		return ret, nil
	})
}
