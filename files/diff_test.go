package files_test

import (
	"context"
	"testing"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/filestesthelp"

	"github.com/stretchr/testify/require"
)

func TestDiff_Validate(t *testing.T) {
	var x *files.Diff
	require.Error(t, x.Validate())
}

func TestCalculateDiff(t *testing.T) {
	ctx := context.Background()
	existing := files.NewSystem[*files.State](map[files.Path]*files.State{
		"nochange":  filestesthelp.NewStateNewFile("hello world"),
		"remove_me": filestesthelp.NewStateNewFile("contents to remove"),
		"update_me": filestesthelp.NewStateNewFile("contents to update"),
		"new_file": {
			FileExistence: files.FileExistenceAbsent,
		},
	})
	desired := files.NewSystem[*files.StateWithChangeReason](map[files.Path]*files.StateWithChangeReason{
		"nochange": {
			State:        *filestesthelp.NewStateNewFile("hello world"),
			ChangeReason: nil,
		},
		"remove_me": {
			State: files.State{
				FileExistence: files.FileExistenceAbsent,
			},
			ChangeReason: nil,
		},
		"update_me": {
			State:        *filestesthelp.NewStateNewFile("updated contents"),
			ChangeReason: nil,
		},
		"new_file": {
			State:        *filestesthelp.NewStateNewFile("new contents"),
			ChangeReason: nil,
		},
	})
	dr, err := files.CalculateDiff(ctx, existing, desired)
	require.NoError(t, err)
	require.Equal(t, 4, len(dr.Paths()))
	require.Equal(t, files.DiffActionNoChange, dr.Get("nochange").Diff.DiffResult.DiffAction)
	require.Equal(t, files.DiffActionDelete, dr.Get("remove_me").Diff.DiffResult.DiffAction)
	require.Equal(t, files.DiffActionUpdate, dr.Get("update_me").Diff.DiffResult.DiffAction)
	require.Equal(t, "updated contents", string(dr.Get("update_me").Diff.NewFileState.Contents))
	require.Equal(t, files.DiffActionCreate, dr.Get("new_file").Diff.DiffResult.DiffAction)
}
