package planner

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/cresta/zapctx"
	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/config/configloader"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/stateloader"
	"github.com/getsyncer/syncer-core/git"
	"github.com/getsyncer/syncer-core/syncer"
	"go.uber.org/zap"
)

type Planner interface {
	Plan(ctx context.Context) (*files.System[*files.DiffWithChangeReason], error)
}

func NewPlanner(registry drift.Registry, configLoader configloader.ConfigLoader, log *zapctx.Logger, stateLoader stateloader.StateLoader, g git.Git, hook Hook) Planner {
	return &plannerImpl{
		registry:     registry,
		configLoader: configLoader,
		log:          log,
		stateLoader:  stateLoader,
		git:          g,
		hook:         hook,
	}
}

type plannerImpl struct {
	registry     drift.Registry
	configLoader configloader.ConfigLoader
	log          *zapctx.Logger
	stateLoader  stateloader.StateLoader
	git          git.Git
	hook         Hook
}

var _ Planner = &plannerImpl{}

func (s *plannerImpl) Plan(ctx context.Context) (*files.System[*files.DiffWithChangeReason], error) {
	s.log.Debug(ctx, "Starting plan")
	rc, err := configloader.ConfigFromFile(ctx, s.configLoader)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	s.log.Debug(ctx, "config pre flatten", zap.Any("config", rc))
	if err := s.configLoader.FlattenChildren(ctx, rc); err != nil {
		return nil, fmt.Errorf("failed to flatten children: %w", err)
	}
	s.log.Debug(ctx, "Loaded config", zap.Any("config", rc))
	printConfigIfDebug(ctx, s.log, rc)
	if err := s.hook.PreSetup(ctx, rc); err != nil {
		return nil, fmt.Errorf("failed to run pre-setup hook: %w", err)
	}
	if err := s.mergeConfigs(ctx, rc); err != nil {
		return nil, fmt.Errorf("failed to merge configs: %w", err)
	}
	rc.Syncs = sortSyncs(rc.Syncs, s.registry)
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	if err := s.loopAndExecute(ctx, rc, wd, loopAndExecuteSetup); err != nil {
		return nil, fmt.Errorf("failed to setup sync: %w", err)
	}
	changes := make([]*files.System[*files.StateWithChangeReason], 0, len(rc.Syncs))
	if err := s.loopAndExecute(ctx, rc, wd, loopAndExecuteRun(&changes)); err != nil {
		return nil, fmt.Errorf("failed to run sync: %w", err)
	}
	s.log.Debug(ctx, "Merging changes", zap.Any("changes", changes))
	finalExpectedState, err := files.SystemMerge(changes...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge changes: %w", err)
	}
	allSyncedFiles, err := stateloader.SyncedFiles(ctx, s.git, s.stateLoader, wd, drift.MagicTrackedString)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracked files: %w", err)
	}
	allSyncedFiles.RemoveAll(finalExpectedState.Paths())
	removals, err := files.ConvertToRemovals(allSyncedFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to removals: %w", err)
	}
	finalExpectedState, err = files.SystemMerge(finalExpectedState, removals)
	if err != nil {
		return nil, fmt.Errorf("failed to merge removals: %w", err)
	}
	allPaths := finalExpectedState.Paths()
	s.log.Debug(ctx, "Loading existing state", zap.Any("paths", allPaths))
	existingState, err := stateloader.LoadAllState(ctx, allPaths, s.stateLoader)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing state: %w", err)
	}
	s.log.Debug(ctx, "Calculating diff", zap.Any("existing", existingState), zap.Any("final", finalExpectedState))
	stateDiff, err := files.CalculateDiff(ctx, existingState, finalExpectedState)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate diff: %w", err)
	}
	return stateDiff, nil
}

func printConfigIfDebug(ctx context.Context, logger *zapctx.Logger, rc *config.Root) {
	if logger.Unwrap(ctx).Level() <= zap.DebugLevel {
		asY, err := rc.AsYaml()
		if err != nil {
			fmt.Printf("Failed to print config as yaml: %v\n", err)
		}
		fmt.Println(asY)
	}
}

func (s *plannerImpl) mergeConfigs(ctx context.Context, rc *config.Root) error {
	for idx := range rc.Syncs {
		r := rc.Syncs[idx]
		s.log.Debug(ctx, "Config before merge", zap.Any("config", r.Config), zap.String("config-as-yaml", valOrErr(r.Config.AsYaml())))
		if err := rc.Syncs[idx].Config.Merge(rc.Config); err != nil {
			return fmt.Errorf("failed to merge run config: %w", err)
		}
		s.log.Debug(ctx, "Config after merge", zap.Any("config", r.Config), zap.String("config-as-yaml", valOrErr(r.Config.AsYaml())))
	}
	return nil
}

func valOrErr(v1 string, err error) string {
	if err != nil {
		return err.Error()
	}
	return v1
}

type loopAndRunLogic func(ctx context.Context, syncer drift.Detector, runData *drift.RunData) error

func (s *plannerImpl) loopAndExecute(ctx context.Context, rc *config.Root, wd string, toRun loopAndRunLogic) error {
	for _, r := range rc.Syncs {
		logic, exists := s.registry.Get(r.Logic)
		s.log.Debug(ctx, "config for this execute", zap.Any("run-config", r.Config))
		if !exists {
			return fmt.Errorf("logic %s not found", r.Logic)
		}
		sr := drift.RunData{
			Registry:              s.registry,
			RootConfig:            rc,
			RunConfig:             r.Config,
			DestinationWorkingDir: wd,
		}
		if err := toRun(ctx, logic, &sr); err != nil {
			return fmt.Errorf("error running %v: %w", logic.Name(), err)
		}
	}
	return nil
}

func loopAndExecuteRun(changes *[]*files.System[*files.StateWithChangeReason]) func(ctx context.Context, syncer drift.Detector, runData *drift.RunData) error {
	return func(ctx context.Context, syncer drift.Detector, runData *drift.RunData) error {
		var runChanges *files.System[*files.StateWithChangeReason]
		var err error
		if runChanges, err = syncer.DetectDrift(ctx, runData); err != nil {
			return fmt.Errorf("error running %v: %w", syncer.Name(), err)
		}
		*changes = append(*changes, runChanges)
		return nil
	}
}

func loopAndExecuteSetup(ctx context.Context, detector drift.Detector, runData *drift.RunData) error {
	if canSetup, ok := detector.(syncer.SetupSyncer); ok {
		if err := canSetup.Setup(ctx, runData); err != nil {
			return fmt.Errorf("error setting up %v: %w", detector.Name(), err)
		}
	}
	return nil
}

func sortSyncs(syncs []config.Syncs, reg drift.Registry) []config.Syncs {
	ret := make([]config.Syncs, 0, len(syncs))
	sort.SliceStable(syncs, func(i, j int) bool {
		iLogic, iExists := reg.Get(syncs[i].Logic)
		jLogic, jExists := reg.Get(syncs[j].Logic)
		if !iExists || !jExists {
			return syncs[i].Logic < syncs[j].Logic
		}
		return iLogic.Priority() > jLogic.Priority()
	})
	ret = append(ret, syncs...)
	return ret
}
