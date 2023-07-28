package syncertesthelper

import (
	"context"

	"github.com/getsyncer/syncer-core/drift"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/files"
)

type TestSyncer struct {
	Ret          *files.System[*files.StateWithChangeReason]
	SyncerName   config.Name
	SyncPriority drift.Priority
}

func NewTestSyncer(ret *files.System[*files.StateWithChangeReason], name config.Name, priority drift.Priority) *TestSyncer {
	return &TestSyncer{
		Ret:          ret,
		SyncerName:   name,
		SyncPriority: priority,
	}
}

func (t *TestSyncer) DetectDrift(_ context.Context, _ *drift.RunData) (*files.System[*files.StateWithChangeReason], error) {
	return t.Ret, nil
}

func (t *TestSyncer) Name() config.Name {
	return t.SyncerName
}

func (t *TestSyncer) Priority() drift.Priority {
	if t.SyncPriority != 0 {
		return t.SyncPriority
	}
	return drift.PriorityNormal
}

var _ drift.Detector = &TestSyncer{}
