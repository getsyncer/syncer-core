package drifttest

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func FileContains(t *testing.T, path string, expected string) {
	t.Helper()
	require.FileExists(t, path)
	contents, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)
	require.Contains(t, string(contents), expected, "file %s does not contain %s", path, expected)
}

func FileDoesNotContain(t *testing.T, path string, expected string) {
	t.Helper()
	require.FileExists(t, path)
	contents, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)
	require.NotContains(t, string(contents), expected, "file %s contains %s", path, expected)
}

func FileContentMatches(t *testing.T, path string, expected string) {
	t.Helper()
	require.FileExists(t, path)
	contents, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)
	require.Equal(t, expected, string(contents), "file %s does not match %s", path, expected)
}

func FileIsYAML(t *testing.T, path string) {
	t.Helper()
	require.FileExists(t, path)
	contents, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)
	require.NoError(t, yaml.Unmarshal(contents, &struct{}{}), "file %s is not valid YAML", path)
}

func OnlyGitChanges(t *testing.T, paths ...string) {
	g, err := git.PlainOpen("")
	require.NoError(t, err)
	wt, err := g.Worktree()
	require.NoError(t, err)

	// Only the files in path are changed
	status, err := wt.Status()
	require.NoError(t, err)
	require.Equal(t, len(paths), len(status))
	for _, path := range paths {
		if status.IsUntracked(path) {
			// Untracked files must exist
			_, err := os.Stat(path)
			require.NoError(t, err)
			continue
		}
		fs := status.File(path)
		require.NotNil(t, fs)
		require.Equal(t, git.Modified, fs.Worktree)
	}
}
