package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type Diff struct {
	OldFileState *State
	NewFileState *State
	DiffResult   DiffResult
}

func (d *Diff) Validate() error {
	return d.DiffResult.Validate()
}

type DiffResult struct {
	DiffAction         DiffAction
	ModeToChangeTo     *os.FileMode
	ContentsToChangeTo []byte
}

func (r *DiffResult) Validate() error {
	if r.DiffAction == DiffActionUnset {
		return fmt.Errorf("diff action must be set")
	}
	if r.DiffAction == DiffActionNoChange {
		if r.ModeToChangeTo != nil {
			return fmt.Errorf("mode must be empty when no change")
		}
		if r.ContentsToChangeTo != nil {
			return fmt.Errorf("contents must be empty when no change")
		}
	}
	if r.DiffAction == DiffActionDelete {
		if r.ModeToChangeTo != nil {
			return fmt.Errorf("mode must be empty when deleting")
		}
		if r.ContentsToChangeTo != nil {
			return fmt.Errorf("contents must be empty when deleting")
		}
	}
	if r.DiffAction == DiffActionCreate {
		if r.ModeToChangeTo == nil {
			return fmt.Errorf("mode must be set when creating")
		}
		if r.ContentsToChangeTo == nil {
			return fmt.Errorf("contents must be set when creating")
		}
	}
	if r.DiffAction == DiffActionUpdate {
		if r.ModeToChangeTo == nil && r.ContentsToChangeTo == nil {
			return fmt.Errorf("mode or contents must be set when updating")
		}
	}
	return nil
}

type DiffAction int

const (
	DiffActionUnset    DiffAction = iota
	DiffActionDelete              // Delete the object
	DiffActionCreate              // Create the object
	DiffActionUpdate              // Update the object
	DiffActionNoChange            // No change to the object
)

func (d DiffAction) String() string {
	switch d {
	case DiffActionUnset:
		return "unset"
	case DiffActionDelete:
		return "delete"
	case DiffActionCreate:
		return "create"
	case DiffActionUpdate:
		return "update"
	case DiffActionNoChange:
		return "no change"
	default:
		panic("unreachable")
	}
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

func ExecuteDiffOnOs(path Path, d *Diff) error {
	if d.DiffResult.DiffAction == DiffActionNoChange {
		return nil
	}
	if d.DiffResult.DiffAction == DiffActionDelete {
		if err := os.Remove(string(path)); err != nil {
			return fmt.Errorf("failed to delete %s: %w", path, err)
		}
		return nil
	}
	if d.DiffResult.DiffAction == DiffActionCreate {
		dirOfFile := filepath.Dir(string(path))
		if err := os.MkdirAll(dirOfFile, 0755); err != nil {
			return fmt.Errorf("failed to mkdir %s: %w", dirOfFile, err)
		}
		if err := os.WriteFile(string(path), d.DiffResult.ContentsToChangeTo, *d.DiffResult.ModeToChangeTo); err != nil {
			return fmt.Errorf("failed to create %s: %w", path, err)
		}
		return nil
	}
	if d.DiffResult.DiffAction == DiffActionUpdate {
		if d.DiffResult.ModeToChangeTo != nil {
			if err := os.Chmod(string(path), *d.DiffResult.ModeToChangeTo); err != nil {
				return fmt.Errorf("failed to chmod %s: %w", path, err)
			}
		}
		if d.DiffResult.ContentsToChangeTo != nil {
			if err := os.WriteFile(string(path), d.DiffResult.ContentsToChangeTo, 0); err != nil {
				return fmt.Errorf("failed to write %s: %w", path, err)
			}
		}
	}
	return nil
}
