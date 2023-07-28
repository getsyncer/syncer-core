package osdiffexecutor

import (
	"github.com/getsyncer/syncer-core/files/diffexecutor"
	"go.uber.org/fx"
)

var Module = fx.Module("osdiffexecutor", fx.Provide(
	fx.Annotate(
		newOsDiffExecutor,
		fx.As(new(diffexecutor.DiffExecutor)))))
