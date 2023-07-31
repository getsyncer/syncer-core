package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) {
	path, err := exec.LookPath("git")
	if path == "" || err != nil {
		t.Skipf("git not found in path: %v", err)
	}
	// Setup an example git repo
	td := t.TempDir()
	runInDir(t, td, "git", "init")
	runInDir(t, td, "git", "config", "--local", "user.email", "noreply@example.com")
	runInDir(t, td, "git", "config", "--local", "user.name", "noreply")
	require.NoError(t, os.WriteFile(filepath.Join(td, "foo.txt"), []byte("foo"), 0600))
	require.NoError(t, os.Mkdir(filepath.Join(td, "bar"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(td, "bar", "baz.txt"), []byte("baz"), 0600))
	runInDir(t, td, "git", "add", ".")
	runInDir(t, td, "git", "commit", "-m", "initial commit")
	// Test
	g := &gitOs{}
	ctx := context.Background()
	files, err := g.ListTrackedFiles(ctx, td)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"foo.txt", "bar/baz.txt"}, files)
	gr, err := g.FindGitRoot(ctx, filepath.Join(td, "bar"))
	require.NoError(t, err)
	require.Equal(t, td, gr)
}

func runInDir(t *testing.T, dir string, args ...string) {
	//nolint:gosec
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	require.NoError(t, cmd.Run())
}
