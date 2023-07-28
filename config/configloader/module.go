package configloader

import "go.uber.org/fx"

var Module = fx.Module("configloader",
	fx.Provide(
		fx.Annotate(
			NewDefaultConfigLoader,
			fx.As(new(ConfigLoader)),
		),
	),
)
