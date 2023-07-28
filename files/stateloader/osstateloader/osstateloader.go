package osstateloader

import (
	"context"
	"fmt"
	"os"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/stateloader"
)

type osLoader struct{}

func newOsLoader() *osLoader {
	return &osLoader{}
}

func (o *osLoader) LoadState(_ context.Context, path files.Path) (*files.State, error) {
	state, err := newStateFromPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load state for %s: %w", path, err)
	}
	return state, nil
}

var _ stateloader.StateLoader = &osLoader{}

func newStateFromPath(path files.Path) (*files.State, error) {
	pathStr := path.String()
	var ret files.State
	var fs os.FileInfo
	var err error
	if fs, err = os.Stat(pathStr); err != nil {
		if os.IsNotExist(err) {
			return &files.State{
				FileExistence: files.FileExistenceAbsent,
			}, nil
		}
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}
	ret.FileExistence = files.FileExistencePresent
	ret.Mode = fs.Mode()
	currentContent, err := os.ReadFile(pathStr)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for new state %s: %w", pathStr, err)
	}
	ret.Contents = currentContent
	return &ret, nil
}
