package fxcli

import (
	"context"

	"go.uber.org/fx"
)

// Module allows invoking CLI style applications from an fx.App
var Module = fx.Module("syncer",
	fx.Provide(
		newFxCLI,
	),
	fx.Invoke(func(_ *fxCLI) {}),
)

// Main is the interface that must be implemented by the main CLI application.
type Main interface {
	// Run is the equivalent of main() in a CLI application.
	Run()
}

type fxCLI struct {
	toRun Main
	sh    fx.Shutdowner
}

func newFxCLI(toRun Main, lc fx.Lifecycle, sh fx.Shutdowner) *fxCLI {
	ret := &fxCLI{
		sh:    sh,
		toRun: toRun,
	}

	lc.Append(fx.Hook{
		OnStart: ret.start,
		OnStop:  ret.stop,
	})

	return ret
}

func (s *fxCLI) start(_ context.Context) error {
	go s.run()
	return nil
}

func (s *fxCLI) stop(_ context.Context) error {
	return nil
}

func (s *fxCLI) run() {
	defer func() {
		if err := s.sh.Shutdown(); err != nil {
			panic(err)
		}
	}()
	s.toRun.Run()
}
