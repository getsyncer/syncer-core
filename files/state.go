package files

import (
	"context"
	"fmt"
	"os"
)

type State struct {
	Mode          os.FileMode
	Contents      []byte
	FileExistence Existence
}

func (f *State) Validate() error {
	if f.FileExistence == FileExistenceUnset {
		return fmt.Errorf("file existence must be set")
	}
	if f.FileExistence == FileExistenceAbsent {
		if f.Contents != nil {
			return fmt.Errorf("contents must be empty when file is absent")
		}
	}
	if f.FileExistence == FileExistencePresent {
		if f.Contents == nil {
			return fmt.Errorf("contents must be set when file is present")
		}
		if f.Mode == 0 {
			return fmt.Errorf("mode must be set when file is present")
		}
	}
	return nil
}

func (f *State) Diff(_ context.Context, newState *State) (*Diff, error) {
	ret := Diff{
		OldFileState: f,
		NewFileState: newState,
	}
	if f.FileExistence == FileExistenceUnset {
		return nil, fmt.Errorf("file existence must be set on receiver")
	}
	if newState.FileExistence == FileExistenceUnset {
		return nil, fmt.Errorf("file existence must be set on argument")
	}
	if f.FileExistence != newState.FileExistence {
		// Deleting an existing file
		if newState.FileExistence == FileExistenceAbsent {
			ret.DiffResult.DiffAction = DiffActionDelete
			return &ret, nil
		}
		// Creating a new file
		if newState.FileExistence == FileExistencePresent && f.FileExistence == FileExistenceAbsent {
			ret.DiffResult = DiffResult{
				DiffAction:         DiffActionCreate,
				ModeToChangeTo:     &newState.Mode,
				ContentsToChangeTo: newState.Contents,
			}
			return &ret, nil
		}
		panic("BUG: unhandled file state") // Should be impossible
	}
	// File should remain deleted
	if f.FileExistence == FileExistenceAbsent && newState.FileExistence == FileExistenceAbsent {
		ret.DiffResult = DiffResult{
			DiffAction: DiffActionNoChange,
		}
		return &ret, nil
	}
	if f.FileExistence != FileExistencePresent || newState.FileExistence != FileExistencePresent {
		panic("BUG: Do not expect present at this point in the code") // Should be impossible
	}
	ret.DiffResult.DiffAction = DiffActionNoChange
	if f.Mode != newState.Mode {
		ret.DiffResult.DiffAction = DiffActionUpdate
		ret.DiffResult.ModeToChangeTo = &newState.Mode
	}
	if string(f.Contents) != string(newState.Contents) {
		ret.DiffResult.DiffAction = DiffActionUpdate
		ret.DiffResult.ContentsToChangeTo = newState.Contents
	}
	return &ret, nil
}

type Existence int

const (
	FileExistenceUnset Existence = iota
	FileExistencePresent
	FileExistenceAbsent
)

func (e Existence) String() string {
	switch e {
	case FileExistenceUnset:
		return "unset"
	case FileExistencePresent:
		return "present"
	case FileExistenceAbsent:
		return "absent"
	default:
		panic("BUG: unknown existence")
	}
}

func SimpleState(files map[string]string) *System[*State] {
	var ret System[*State]
	for path, contents := range files {
		if err := ret.Add(Path(path), &State{
			Mode:          0644,
			Contents:      []byte(contents),
			FileExistence: FileExistencePresent,
		}); err != nil {
			panic(err)
		}
	}
	return &ret
}
