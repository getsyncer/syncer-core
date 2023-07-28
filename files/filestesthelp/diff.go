package filestesthelp

import (
	"os"

	"github.com/getsyncer/syncer-core/files"
)

func NewDiffNewFile(contents string) *files.Diff {
	m := os.FileMode(0644)
	return &files.Diff{
		NewFileState: &files.State{
			FileExistence: files.FileExistencePresent,
			Contents:      []byte(contents),
			Mode:          m,
		},
		DiffResult: files.DiffResult{
			DiffAction:         files.DiffActionCreate,
			ContentsToChangeTo: []byte(contents),
			ModeToChangeTo:     &m,
		},
	}
}

func NewStateNewFile(contents string) *files.State {
	m := os.FileMode(0644)
	return &files.State{
		FileExistence: files.FileExistencePresent,
		Contents:      []byte(contents),
		Mode:          m,
	}
}

func NewDiffNoChange(state *files.State) *files.Diff {
	return &files.Diff{
		OldFileState: state,
		NewFileState: state,
		DiffResult: files.DiffResult{
			DiffAction: files.DiffActionNoChange,
		},
	}
}

func NewDiffChangeContent(oldState *files.State, newContents string) *files.Diff {
	return &files.Diff{
		OldFileState: oldState,
		NewFileState: &files.State{
			FileExistence: files.FileExistencePresent,
			Contents:      []byte(newContents),
			Mode:          oldState.Mode,
		},
		DiffResult: files.DiffResult{
			DiffAction:         files.DiffActionUpdate,
			ContentsToChangeTo: []byte(newContents),
		},
	}
}

func NewDiffDeleteFile(state *files.State) *files.Diff {
	return &files.Diff{
		OldFileState: state,
		NewFileState: &files.State{
			FileExistence: files.FileExistenceAbsent,
		},
		DiffResult: files.DiffResult{
			DiffAction: files.DiffActionDelete,
		},
	}
}
