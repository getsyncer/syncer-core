package syncer

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/fileprinter"
	"github.com/getsyncer/syncer-core/fxcli"
	"go.uber.org/fx"
)

const (
	FxTagSyncers  = `group:"syncers"`
	FxTagChildren = `group:"childrensource"`
)

var Module = fx.Module("syncer",
	fx.Provide(
		fx.Annotate(
			NewDefaultConfigLoader,
			fx.As(new(ConfigLoader)),
		),
		fx.Annotate(
			NewRegistry,
			fx.As(new(Registry)),
			fx.ParamTags(FxTagSyncers),
		),
		fx.Annotate(
			NewChildrenRegistry,
			fx.As(new(ChildrenRegistry)),
			fx.ParamTags(FxTagChildren),
		),
		fx.Annotate(
			NewPlanner,
			fx.As(new(Planner)),
		),
		fx.Annotate(
			NewApplier,
			fx.As(new(Applier)),
		),
	),
)

type globalFxRegistry struct {
	options []fx.Option
	mu      sync.Mutex
}

var globalFxRegistryInstance = &globalFxRegistry{}

func (g *globalFxRegistry) Register(opt fx.Option) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.options = append(g.options, opt)
}

func (g *globalFxRegistry) Get() []fx.Option {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.options
}

func FxRegister(opt fx.Option) {
	globalFxRegistryInstance.Register(opt)
}

var ExecuteCliModule = fx.Module(
	"main",
	fxcli.Module,
	fx.Provide(
		fx.Annotate(
			NewFxCli,
			fx.As(new(fxcli.Main)),
		),
	),
)

type FxCli struct {
	planner Planner
	applier Applier
	printer fileprinter.Printer
}

func NewFxCli(planner Planner, applier Applier, printer fileprinter.Printer) *FxCli {
	return &FxCli{planner: planner, applier: applier, printer: printer}
}

func (f *FxCli) Run() {
	ctx := context.Background()
	cmd := os.Getenv("SYNCER_EXEC_CMD")
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
			fmt.Println("Error: ", err)
			return
		}
		if err := f.printer.PrettyPrintDiffs(os.Stdout, diffs); err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if cmd == plan {
			if os.Getenv("SYNCER_EXIT_CODE_ON_DIFF") == "true" {
				if files.IncludesChanges(diffs) {
					os.Exit(1)
				}
			}
		}
		if cmd == apply {
			if err := f.applier.Apply(ctx, diffs); err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
		return
	default:
		fmt.Println("Unknown command: ", cmd)
	}
}
