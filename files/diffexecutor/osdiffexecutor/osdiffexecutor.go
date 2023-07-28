package osdiffexecutor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/diffexecutor"
)

type osDiffExecutor struct{}

func newOsDiffExecutor() *osDiffExecutor {
	return &osDiffExecutor{}
}

func (o *osDiffExecutor) ExecuteDiff(_ context.Context, path files.Path, d *files.Diff) error {
	if err := executeDiffOnOs(path, d); err != nil {
		return fmt.Errorf("failed to execute diff for %s: %w", path, err)
	}
	return nil
}

var _ diffexecutor.DiffExecutor = &osDiffExecutor{}

func executeDiffOnOs(path files.Path, d *files.Diff) error {
	if d.DiffResult.DiffAction == files.DiffActionNoChange {
		return nil
	}
	if d.DiffResult.DiffAction == files.DiffActionDelete {
		if err := os.Remove(string(path)); err != nil {
			return fmt.Errorf("failed to delete %s: %w", path, err)
		}
		return nil
	}
	if d.DiffResult.DiffAction == files.DiffActionCreate {
		dirOfFile := filepath.Dir(string(path))
		if err := os.MkdirAll(dirOfFile, 0755); err != nil {
			return fmt.Errorf("failed to mkdir %s: %w", dirOfFile, err)
		}
		if err := os.WriteFile(string(path), d.DiffResult.ContentsToChangeTo, *d.DiffResult.ModeToChangeTo); err != nil {
			return fmt.Errorf("failed to create %s: %w", path, err)
		}
		return nil
	}
	if d.DiffResult.DiffAction == files.DiffActionUpdate {
		if d.DiffResult.ModeToChangeTo != nil {
			if err := os.Chmod(string(path), *d.DiffResult.ModeToChangeTo); err != nil {
				return fmt.Errorf("failed to chmod %s: %w", path, err)
			}
		}
		if d.DiffResult.ContentsToChangeTo != nil {
			if err := os.WriteFile(string(path), d.DiffResult.ContentsToChangeTo, 0); err != nil {
				return fmt.Errorf("failed to write %s: %w", path, err)
			}
		}
	}
	return nil
}
