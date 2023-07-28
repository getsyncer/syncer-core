package osdiffexecutor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/getsyncer/syncer-core/files/filestesthelp"

	"github.com/getsyncer/syncer-core/files"

	"github.com/stretchr/testify/require"
)

func TestOsDiffExecutor_ExecuteDiff(t *testing.T) {
	x := newOsDiffExecutor()
	td := t.TempDir()
	ctx := context.Background()
	newContent := filestesthelp.NewDiffNewFile("hello world")
	filePath := files.Path(filepath.Join(td, "a"))
	require.NoError(t, x.ExecuteDiff(ctx, filePath, newContent))
	cont, err := os.ReadFile(filepath.Join(td, "a"))
	require.NoError(t, err)
	require.Equal(t, "hello world", string(cont))

	require.NoError(t, x.ExecuteDiff(ctx, filePath, filestesthelp.NewDiffChangeContent(newContent.NewFileState, "hello world 2")))
	cont, err = os.ReadFile(filepath.Join(td, "a"))
	require.NoError(t, err)
	require.Equal(t, "hello world 2", string(cont))

	require.NoError(t, x.ExecuteDiff(ctx, filePath, filestesthelp.NewDiffDeleteFile(newContent.NewFileState)))
	_, err = os.Stat(filepath.Join(td, "a"))
	require.True(t, os.IsNotExist(err))
}
