package testdiffexecutor

import (
	"github.com/getsyncer/syncer-core/files/diffexecutor"
	"go.uber.org/fx"
)

var Module = fx.Module("testdiffexecutor", fx.Provide(
	fx.Annotate(
		newTest,
		fx.As(new(diffexecutor.DiffExecutor)))))
