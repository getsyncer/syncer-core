package files_test

import (
	"context"
	"os"
	"testing"

	"github.com/getsyncer/syncer-core/files"

	"github.com/stretchr/testify/require"
)

func TestState_Diff(t *testing.T) {
	expectErrorDiff := func(t *testing.T, oldState, newState *files.State) func(t *testing.T) {
		return func(t *testing.T) {
			_, err := oldState.Diff(context.Background(), newState)
			require.Error(t, err)
		}
	}
	expectNoErrorDiff := func(t *testing.T, oldState, newState *files.State, expectedDiff *files.Diff) func(t *testing.T) {
		return func(t *testing.T) {
			diffRes, err := oldState.Diff(context.Background(), newState)
			require.NoError(t, err)
			require.Equal(t, expectedDiff, diffRes)
		}
	}
	t.Run("file existence must be set", expectErrorDiff(t, &files.State{}, &files.State{}))
	t.Run("file existence must be set on receiver", expectErrorDiff(t, &files.State{FileExistence: files.FileExistenceUnset}, &files.State{FileExistence: files.FileExistencePresent}))
	t.Run("file existence must be set on argument", expectErrorDiff(t, &files.State{FileExistence: files.FileExistencePresent}, &files.State{FileExistence: files.FileExistenceUnset}))

	t.Run("no change if both absent", expectNoErrorDiff(t, &files.State{FileExistence: files.FileExistenceAbsent}, &files.State{FileExistence: files.FileExistenceAbsent}, &files.Diff{
		OldFileState: &files.State{FileExistence: files.FileExistenceAbsent},
		NewFileState: &files.State{FileExistence: files.FileExistenceAbsent},
		DiffResult: files.DiffResult{
			DiffAction: files.DiffActionNoChange,
		},
	}))
	t.Run("no change if both present same content", expectNoErrorDiff(t, &files.State{FileExistence: files.FileExistencePresent}, &files.State{FileExistence: files.FileExistencePresent}, &files.Diff{
		OldFileState: &files.State{FileExistence: files.FileExistencePresent},
		NewFileState: &files.State{FileExistence: files.FileExistencePresent},
		DiffResult: files.DiffResult{
			DiffAction: files.DiffActionNoChange,
		},
	}))
	t.Run("delete if old present, new absent", expectNoErrorDiff(t, &files.State{FileExistence: files.FileExistencePresent}, &files.State{FileExistence: files.FileExistenceAbsent}, &files.Diff{
		OldFileState: &files.State{FileExistence: files.FileExistencePresent},
		NewFileState: &files.State{FileExistence: files.FileExistenceAbsent},
		DiffResult: files.DiffResult{
			DiffAction: files.DiffActionDelete,
		},
	}))
	m := os.FileMode(0644)
	t.Run("create if old absent, new present", expectNoErrorDiff(t, &files.State{FileExistence: files.FileExistenceAbsent}, &files.State{FileExistence: files.FileExistencePresent, Mode: m, Contents: []byte("hello world")}, &files.Diff{
		OldFileState: &files.State{FileExistence: files.FileExistenceAbsent},
		NewFileState: &files.State{
			FileExistence: files.FileExistencePresent,
			Mode:          m,
			Contents:      []byte("hello world"),
		},
		DiffResult: files.DiffResult{
			DiffAction:         files.DiffActionCreate,
			ModeToChangeTo:     &m,
			ContentsToChangeTo: []byte("hello world"),
		},
	}))
	t.Run("update if old present, new present", expectNoErrorDiff(t, &files.State{FileExistence: files.FileExistencePresent, Mode: m, Contents: []byte("hello world")}, &files.State{FileExistence: files.FileExistencePresent, Mode: m, Contents: []byte("hello world2")}, &files.Diff{
		OldFileState: &files.State{
			FileExistence: files.FileExistencePresent,
			Mode:          m,
			Contents:      []byte("hello world"),
		},
		NewFileState: &files.State{
			FileExistence: files.FileExistencePresent,
			Mode:          m,
			Contents:      []byte("hello world2"),
		},
		DiffResult: files.DiffResult{
			DiffAction:         files.DiffActionUpdate,
			ContentsToChangeTo: []byte("hello world2"),
		},
	}))
}
