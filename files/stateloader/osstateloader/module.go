package osstateloader

import (
	"github.com/getsyncer/syncer-core/files/stateloader"
	"go.uber.org/fx"
)

var Module = fx.Module("osstateloader", fx.Provide(
	fx.Annotate(
		newOsLoader,
		fx.As(new(stateloader.StateLoader)))))
