package staticfile

import (
	"github.com/getsyncer/syncer-core/syncer"
	"go.uber.org/fx"
)

func init() {
	syncer.FxRegister(Module)
}

var Module = fx.Module("staticfile",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(syncer.DriftSyncer)),
			fx.ResultTags(`group:"syncers"`),
		),
	),
)
