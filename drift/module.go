package drift

import (
	"go.uber.org/fx"
)

const (
	FxTagSyncers = `group:"syncers"`
)

var Module = fx.Module("syncer",
	fx.Provide(
		fx.Annotate(
			NewRegistry,
			fx.As(new(Registry)),
			fx.ParamTags(FxTagSyncers),
		),
	),
)
