package staticfile

import (
	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files/stateloader"
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

func NewCustomModule[T CreateFilesystem](name config.Name, priority drift.Priority) fx.Option {
	return fx.Module(name.String(),
		fx.Provide(
			fx.Annotate(
				func(loader stateloader.StateLoader) *Syncer[T] {
					return NewWithCustomLogic[T](name, priority, loader)
				},
				fx.As(new(drift.Detector)),
				fx.ResultTags(`group:"syncers"`),
			),
		),
	)
}
