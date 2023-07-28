package drift

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/files"
)

type Priority int

const (
	PriorityLowest  = Priority(100)
	PriorityLow     = Priority(200)
	PriorityNormal  = Priority(300)
	PriorityHigh    = Priority(400)
	PriorityHighest = Priority(500)
)

func (p Priority) String() string {
	switch p {
	case PriorityLowest:
		return "lowest"
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityHighest:
		return "highest"
	default:
		return fmt.Sprintf("unknown(%d)", p)
	}
}

type Detector interface {
	DetectDrift(ctx context.Context, runData *RunData) (*files.System[*files.StateWithChangeReason], error)
	Name() config.Name
	Priority() Priority
}
