package stateloader_test

import (
	"context"
	"testing"

	"github.com/getsyncer/syncer-core/git/testgit"

	"github.com/getsyncer/syncer-core/files/stateloader"
	"github.com/stretchr/testify/require"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/stateloader/teststateloader"
)

func TestLoadAllState(t *testing.T) {
	ctx := context.Background()
	ts := teststateloader.NewTestStateloaderFromFiles(map[string]string{
		"readme":  "hello world",
		"main.go": "func main() {}",
	})
	paths := []files.Path{"readme"}
	fsState, err := stateloader.LoadAllState(ctx, paths, ts)
	require.NoError(t, err)
	require.Equal(t, 1, len(fsState.Paths()))
	require.Equal(t, "hello world", string(fsState.Get("readme").Contents))

}

func TestSyncedFiles(t *testing.T) {
	ctx := context.Background()
	ts := teststateloader.NewTestStateloaderFromFiles(map[string]string{
		"readme":  "hello world AUTOGEN FILE",
		"main.go": "func main() {}",
	})
	g := testgit.New("/path", []string{"readme", "main.go"})
	fsState, err := stateloader.SyncedFiles(ctx, g, ts, "/path", "AUTOGEN FILE")
	require.NoError(t, err)
	require.Equal(t, 1, len(fsState.Paths()))
	require.Equal(t, "hello world AUTOGEN FILE", string(fsState.Get("readme").Contents))
}
