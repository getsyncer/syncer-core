package templatefiles

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files"
)

type PostGenProcessor interface {
	PostGenProcess(ctx context.Context, fs *files.System[*files.StateWithChangeReason], runData *drift.RunData) error
}

type PostGenProcessorList []PostGenProcessor

func (p PostGenProcessorList) PostGenProcess(ctx context.Context, fs *files.System[*files.StateWithChangeReason], runData *drift.RunData) error {
	for _, v := range p {
		if err := v.PostGenProcess(ctx, fs, runData); err != nil {
			return fmt.Errorf("unable to post gen process: %w", err)
		}
	}
	return nil
}
