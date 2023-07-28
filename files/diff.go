package files

import (
	"context"
	"fmt"
)

type Diff struct {
	OldFileState *State
	NewFileState *State
	DiffResult   DiffResult
}

func (d *Diff) Validate() error {
	if d == nil {
		return fmt.Errorf("diff is nil")
	}
	return d.DiffResult.Validate()
}

func CalculateDiff(ctx context.Context, existing *System[*State], desired *System[*StateWithChangeReason]) (*System[*DiffWithChangeReason], error) {
	var ret System[*DiffWithChangeReason]
	existingPaths := existing.Paths()
	desiredPaths := desired.Paths()
	seenPaths := map[Path]struct{}{}
	for _, path := range desiredPaths {
		seenPaths[path] = struct{}{}
		if !existing.IsTracked(path) {
			return nil, fmt.Errorf("path %q is not tracked", path)
		}
		asExisting := existing.Get(path)
		asDesired := desired.Get(path)
		diff, err := asExisting.Diff(ctx, &asDesired.State)
		if err != nil {
			return nil, fmt.Errorf("cannot calculate diff for %q: %w", path, err)
		}
		toAdd := &DiffWithChangeReason{
			ChangeReason: asDesired.ChangeReason,
			Diff:         diff,
		}
		if err := ret.Add(path, toAdd); err != nil {
			return nil, fmt.Errorf("cannot add diff for %q: %w", path, err)
		}
	}
	for _, e := range existingPaths {
		if _, ok := seenPaths[e]; !ok {
			return nil, fmt.Errorf("path %q is not desired but was in existing state", e)
		}
	}
	return &ret, nil
}

func IncludesChanges(diffs *System[*DiffWithChangeReason]) bool {
	for _, path := range diffs.Paths() {
		f := diffs.Get(path)
		if f.Diff.DiffResult.DiffAction != DiffActionNoChange {
			return true
		}
	}
	return false
}
