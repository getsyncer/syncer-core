package stateloader

import (
	"context"
	"fmt"
	"strings"

	"github.com/getsyncer/syncer-core/config/configloader"

	"github.com/getsyncer/syncer-core/git"

	"github.com/getsyncer/syncer-core/files"
)

type StateLoader interface {
	LoadState(ctx context.Context, path files.Path) (*files.State, error)
}

func LoadAllState(ctx context.Context, paths []files.Path, loader StateLoader) (*files.System[*files.State], error) {
	var ret files.System[*files.State]
	for _, path := range paths {
		state, err := loader.LoadState(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to load state for %s: %w", path, err)
		}
		if err := ret.Add(path, state); err != nil {
			return nil, fmt.Errorf("failed to add state for %s: %w", path, err)
		}
	}
	return &ret, nil
}

func containsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func SyncedFiles(ctx context.Context, g git.Git, loader StateLoader, loc string, syncFlag string) (*files.System[*files.State], error) {
	trackedFiles, err := g.ListTrackedFiles(ctx, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to list git tracked files: %w", err)
	}
	var ret files.System[*files.State]
	ignoredFiles := configloader.DefaultLocations()
	ignoredFiles = append(ignoredFiles, ".syncer/sync.go")
	for _, f := range trackedFiles {
		f := files.Path(f).Clean().String()
		if containsString(ignoredFiles, f) {
			continue
		}
		s, err := loader.LoadState(ctx, files.Path(f))
		if err != nil {
			return nil, fmt.Errorf("failed to load state for %s: %w", f, err)
		}
		if strings.Contains(string(s.Contents), syncFlag) {
			if err := ret.Add(files.Path(f), s); err != nil {
				return nil, fmt.Errorf("failed to add state for %s: %w", f, err)
			}
		}
	}
	return &ret, nil
}
