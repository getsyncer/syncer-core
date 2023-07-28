package testgit

import (
	"github.com/getsyncer/syncer-core/git"
	"go.uber.org/fx"
)

func GitModule(location string, trackedFiles []string) fx.Option {
	constructor := func() *TestGit {
		return New(location, trackedFiles)
	}
	return fx.Module("testgit", fx.Provide(fx.Annotate(constructor, fx.As(new(git.Git)))))
}
