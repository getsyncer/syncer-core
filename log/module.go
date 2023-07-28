package log

import (
	"go.uber.org/fx"
)

const FxTagLogDynamicFields = `group:"dynamic_fields"`

var ModuleCore = fx.Module("logger",
	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags("", FxTagLogDynamicFields),
		),
	),
)

var ModuleProd = fx.Module("logger",
	fx.Provide(
		ZapLoggerFromConfig,
		ZapLoggerConfigFromEnv,
	),
)
