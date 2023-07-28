package testgit

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/getsyncer/syncer-core/files"

	"github.com/getsyncer/syncer-core/git"
)

type TestGit struct {
	Location     string
	TrackedFiles []string
}

func NewFromFilesystem(location string, paths []files.Path) *TestGit {
	trackedFiles := make([]string, len(paths))
	for i, path := range paths {
		trackedFiles[i] = string(path)
	}
	return New(location, trackedFiles)
}

func New(location string, trackedFiles []string) *TestGit {
	return &TestGit{
		Location:     location,
		TrackedFiles: trackedFiles,
	}
}

func (t *TestGit) FindGitRoot(_ context.Context, loc string) (string, error) {
	loc = filepath.Clean(loc)
	if strings.Contains(loc, t.Location) {
		return t.Location, nil
	}
	return "", fmt.Errorf("not found")
}

func (t *TestGit) ListTrackedFiles(ctx context.Context, loc string) ([]string, error) {
	_, err := t.FindGitRoot(ctx, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to find git root: %w", err)
	}
	ret := make([]string, len(t.TrackedFiles))
	copy(ret, t.TrackedFiles)
	return ret, nil
}

var _ git.Git = (*TestGit)(nil)
