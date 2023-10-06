package planner

import (
	"context"

	"github.com/getsyncer/syncer-core/files"
)

type contextKey int

const (
	currentChangesKey contextKey = iota
)

func WithCurrentChanges(ctx context.Context, changes []*files.System[*files.StateWithChangeReason]) context.Context {
	return context.WithValue(ctx, currentChangesKey, changes)
}

func GetCurrentChanges(ctx context.Context) []*files.System[*files.StateWithChangeReason] {
	toRet := ctx.Value(currentChangesKey)
	if toRet == nil {
		return nil
	}
	return toRet.([]*files.System[*files.StateWithChangeReason])
}
