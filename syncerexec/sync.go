package syncerexec

import (
	"github.com/getsyncer/syncer-core/config/configloader"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files/diffexecutor/osdiffexecutor"
	"github.com/getsyncer/syncer-core/files/fileprinter/consoleprinter"
	"github.com/getsyncer/syncer-core/files/stateloader/osstateloader"
	"github.com/getsyncer/syncer-core/fxcli"
	"github.com/getsyncer/syncer-core/fxregistry"
	"github.com/getsyncer/syncer-core/git"
	"github.com/getsyncer/syncer-core/log"
	"github.com/getsyncer/syncer-core/osstub"
	"github.com/getsyncer/syncer-core/syncer/applier"
	"github.com/getsyncer/syncer-core/syncer/childrenregistry"
	"github.com/getsyncer/syncer-core/syncer/planner"
	"go.uber.org/fx"
)

func DefaultFxOptions() fx.Option {
	var allOpts []fx.Option
	allOpts = append(allOpts, fx.WithLogger(log.NewFxLogger), log.ModuleCore, consoleprinter.Module, fxcli.Module, git.Module, drift.Module, osstateloader.Module, configloader.Module, childrenregistry.Module, osdiffexecutor.Module, planner.Module, applier.Module)
	allOpts = append(allOpts, fxregistry.Get()...)
	return fx.Options(allOpts...)
}

func FromCli(opts ...fx.Option) {
	var allOpts []fx.Option
	// TODO: Fix
	allOpts = append(allOpts, opts...)
	allOpts = append(allOpts, log.ModuleProd, Module, osstub.Module)
	fx.New(allOpts...).Run()
}
