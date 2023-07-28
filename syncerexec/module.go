package syncerexec

import (
	"github.com/getsyncer/syncer-core/fxcli"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"main",
	fx.Provide(
		fx.Annotate(
			NewRunSyncer,
			fx.As(new(fxcli.Main)),
		),
	),
)

var ModuleTest = fx.Module(
	"main",
	fx.Provide(
		NewRunSyncer,
	),
)
