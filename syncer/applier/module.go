package applier

import (
	"go.uber.org/fx"
)

var Module = fx.Module("planner",
	fx.Provide(
		fx.Annotate(
			NewApplier,
			fx.As(new(Applier)),
		),
	),
)
