package planner

import (
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

var Module = fx.Module("planner",
	fx.Provide(
		fx.Annotate(
			NewPlanner,
			fx.As(new(Planner)),
		),
		fx.Annotate(
			NewMultiHook,
			fx.ParamTags(FxTagHooks),
			fx.As(new(Hook)),
		),
	),
)
