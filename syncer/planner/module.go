package planner

import (
	"github.com/cresta/zapctx"
	"github.com/getsyncer/syncer-core/config/configloader"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files/stateloader"
	"github.com/getsyncer/syncer-core/git"
	"go.uber.org/fx"
)

const (
	FxTagHooks = `group:"hooks"`
)

func HookModule(name string, constructor interface{}) fx.Option {
	return fx.Module(
		name,
		fx.Provide(
			fx.Annotate(
				constructor,
				fx.ResultTags(FxTagHooks),
				fx.As(new(Hook)),
			),
		),
	)
}

type newPlannerConfig struct {
	fx.In
	Registry     drift.Registry
	ConfigLoader configloader.ConfigLoader
	Log          *zapctx.Logger
	StateLoader  stateloader.StateLoader
	G            git.Git
	Hook         Hook
	Options      []Option `group:"planoption"`
}

func newPlannerFromConfig(cfg newPlannerConfig) Planner {
	return NewPlanner(cfg.Registry, cfg.ConfigLoader, cfg.Log, cfg.StateLoader, cfg.G, cfg.Hook, cfg.Options)
}

var Module = fx.Module("planner",
	fx.Provide(
		fx.Annotate(
			newPlannerFromConfig,
			fx.As(new(Planner)),
		),
		fx.Annotate(
			NewMultiHook,
			fx.ParamTags(FxTagHooks),
			fx.As(new(Hook)),
		),
	),
)
