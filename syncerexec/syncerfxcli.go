package syncerexec

import (
	"context"
	"os"

	"github.com/getsyncer/syncer-core/osstub"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/fileprinter"
	"github.com/getsyncer/syncer-core/fxcli"
	"github.com/getsyncer/syncer-core/syncer/applier"
	"github.com/getsyncer/syncer-core/syncer/planner"
)

type RunSyncer struct {
	planner planner.Planner
	applier applier.Applier
	printer fileprinter.Printer
	osStub  osstub.OsStub
}

func NewRunSyncer(planner planner.Planner, applier applier.Applier, printer fileprinter.Printer, osStub osstub.OsStub) *RunSyncer {
	return &RunSyncer{planner: planner, applier: applier, printer: printer, osStub: osStub}
}

var _ fxcli.Main = (*RunSyncer)(nil)

func (f *RunSyncer) Run() {
	ctx := context.Background()
	cmd := f.osStub.Getenv("SYNCER_EXEC_CMD")
	const (
		plan  = "plan"
		apply = "apply"
	)
	if cmd == "" {
		cmd = plan
	}
	switch cmd {
	case plan:
		fallthrough
	case apply:
		diffs, err := f.planner.Plan(ctx)
		if err != nil {
			f.osStub.Println("Error: ", err)
			f.osStub.Exit(1)
			return
		}
		if err := f.printer.PrettyPrintDiffs(os.Stdout, diffs); err != nil {
			f.osStub.Println("Error: ", err)
			f.osStub.Exit(1)
			return
		}
		if cmd == plan {
			if f.osStub.Getenv("SYNCER_EXIT_CODE_ON_DIFF") == "true" {
				if files.IncludesChanges(diffs) {
					f.osStub.Exit(1)
					return
				}
			}
		}
		if cmd == apply {
			if err := f.applier.Apply(ctx, diffs); err != nil {
				f.osStub.Println("Error: ", err)
				f.osStub.Exit(1)
				return
			}
		}
		f.osStub.Exit(0)
		return
	default:
		f.osStub.Println("Unknown command: ", cmd)
		f.osStub.Exit(1)
	}
}
