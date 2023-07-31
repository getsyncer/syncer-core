package osstateloader

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/getsyncer/syncer-core/files"

	"github.com/stretchr/testify/require"
)

func Test_osLoader_LoadState(t *testing.T) {
	var x osLoader
	td := t.TempDir()
	ctx := context.Background()
	filePath1 := filepath.Join(td, "file1")
	filePath2 := filepath.Join(td, "file2")
	require.NoError(t, os.WriteFile(filePath1, []byte("hello"), 0600))
	s1, err := x.LoadState(ctx, files.Path(filePath1))
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), s1.Contents)
	s2, err := x.LoadState(ctx, files.Path(filePath2))
	require.NoError(t, err)
	require.Equal(t, files.FileExistenceAbsent, s2.FileExistence)
}
