package teststateloader

import (
	"context"
	"fmt"
	"os"

	"github.com/getsyncer/syncer-core/files"
)

type TestStateloader struct {
	FileSystem *files.System[*files.State]
	ErrorPaths []files.Path
}

func NewTestStateloaderFromFiles(filesToMake map[string]string) *TestStateloader {
	in := make(map[files.Path]*files.State, len(filesToMake))
	for path, contents := range filesToMake {
		in[files.Path(path)] = &files.State{
			Contents:      []byte(contents),
			FileExistence: files.FileExistencePresent,
			Mode:          os.FileMode(0644),
		}
	}
	return &TestStateloader{
		FileSystem: files.NewSystem(in),
	}
}

func (l *TestStateloader) LoadState(_ context.Context, path files.Path) (*files.State, error) {
	for _, errorPath := range l.ErrorPaths {
		if path == errorPath {
			return nil, fmt.Errorf("error loading state for %s", path)
		}
	}
	if !l.FileSystem.IsTracked(path) {
		return &files.State{
			FileExistence: files.FileExistenceAbsent,
		}, nil
	}
	return l.FileSystem.Get(path), nil
}
