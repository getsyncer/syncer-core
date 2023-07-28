package diffexecutor

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/files"
)

type DiffExecutor interface {
	ExecuteDiff(ctx context.Context, path files.Path, d *files.Diff) error
}

type ApplyResult struct{}

func ExecuteAllDiffs(ctx context.Context, s *files.System[*files.DiffWithChangeReason], executor DiffExecutor) error {
	for _, path := range s.Paths() {
		if err := executor.ExecuteDiff(ctx, path, s.Get(path).Diff); err != nil {
			return fmt.Errorf("failed to execute diff for %s: %w", path, err)
		}
	}
	return nil
}
