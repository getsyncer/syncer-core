package applier

import (
	"context"
	"fmt"

	"github.com/cresta/zapctx"
	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/diffexecutor"
	"go.uber.org/zap"
)

type Applier interface {
	Apply(ctx context.Context, stateDiff *files.System[*files.DiffWithChangeReason]) error
}

func NewApplier(log *zapctx.Logger, diffExecutor diffexecutor.DiffExecutor) Applier {
	return &applyImpl{
		log:          log,
		diffExecutor: diffExecutor,
	}
}

type applyImpl struct {
	log          *zapctx.Logger
	diffExecutor diffexecutor.DiffExecutor
}

func (s *applyImpl) Apply(ctx context.Context, stateDiff *files.System[*files.DiffWithChangeReason]) error {
	s.log.Debug(ctx, "Executing diff", zap.Any("diff", stateDiff))
	if err := diffexecutor.ExecuteAllDiffs(ctx, stateDiff, s.diffExecutor); err != nil {
		return fmt.Errorf("failed to execute diff: %w", err)
	}
	return nil
}
