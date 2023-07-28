package teststateloader

import (
	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/stateloader"
	"go.uber.org/fx"
)

func NewModule(fs *files.System[*files.State]) fx.Option {
	constructor := func() *TestStateloader {
		return &TestStateloader{
			FileSystem: fs,
		}
	}
	return fx.Options(
		fx.Provide(fx.Annotate(constructor, fx.As(new(stateloader.StateLoader)))),
	)
}
