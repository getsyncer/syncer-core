package files

import (
	"fmt"
	"os"
)

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
