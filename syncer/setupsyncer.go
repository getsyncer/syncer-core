package syncer

import (
	"context"

	"github.com/getsyncer/syncer-core/drift"
)

type SetupSyncer interface {
	Setup(ctx context.Context, runData *drift.RunData) error
}

type SetupSyncerFunc func(ctx context.Context, runData *drift.RunData) error

func (s SetupSyncerFunc) Setup(ctx context.Context, runData *drift.RunData) error {
	return s(ctx, runData)
}

var _ SetupSyncer = SetupSyncerFunc(nil)

type MultiSetupSyncer []SetupSyncer

func (m MultiSetupSyncer) Setup(ctx context.Context, runData *drift.RunData) error {
	for _, s := range m {
		if err := s.Setup(ctx, runData); err != nil {
			return err
		}
	}
	return nil
}

var _ SetupSyncer = MultiSetupSyncer(nil)
