package planner

import (
	"context"

	"github.com/getsyncer/syncer-core/config"
)

type Hook interface {
	PreSetup(ctx context.Context, rc *config.Root) error
}

type HookImpl struct {
	PreSetupFunc func(ctx context.Context, rc *config.Root) error
}

func (h HookImpl) PreSetup(ctx context.Context, rc *config.Root) error {
	if h.PreSetupFunc == nil {
		return nil
	}
	return h.PreSetupFunc(ctx, rc)
}

type MultiHook struct {
	Hooks []Hook
}

func (m MultiHook) PreSetup(ctx context.Context, rc *config.Root) error {
	for _, h := range m.Hooks {
		if err := h.PreSetup(ctx, rc); err != nil {
			return err
		}
	}
	return nil
}

func NewMultiHook(hooks []Hook) MultiHook {
	return MultiHook{
		Hooks: hooks,
	}
}
