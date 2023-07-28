package staticfile

import (
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/fxregistry"
	"go.uber.org/fx"
)

func init() {
	fxregistry.Register(Module)
}

var Module = fx.Module("staticfile",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(drift.Detector)),
			fx.ResultTags(`group:"syncers"`),
		),
	),
)
