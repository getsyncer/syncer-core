package planner

import (
	"go.uber.org/fx"
)

var Module = fx.Module("planner",
	fx.Provide(
		fx.Annotate(
			NewPlanner,
			fx.As(new(Planner)),
		),
	),
)
