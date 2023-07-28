package staticfile

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/drift"

	"github.com/getsyncer/syncer-core/config"

	"github.com/getsyncer/syncer-core/files"
)

type Config struct {
	// Filename is the name of the file to create
	Filename string `yaml:"filename"`
	// Content is the content of the file to create
	Content string `yaml:"content"`
}

// New returns a new staticfile syncer
func New() *Syncer {
	return &Syncer{}
}

type Syncer struct {
}

// DetectDrift returns a system containing the file specified in the config
func (f *Syncer) DetectDrift(_ context.Context, runData *drift.RunData) (*files.System[*files.StateWithChangeReason], error) {
	var cfg Config
	if err := runData.RunConfig.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staticfile config: %w", err)
	}
	if cfg.Filename == "" {
		return nil, fmt.Errorf("filename is required")
	}
	var ret files.System[*files.StateWithChangeReason]
	path := cfg.Filename
	if err := ret.Add(files.Path(cfg.Filename), &files.StateWithChangeReason{
		ChangeReason: &files.ChangeReason{
			Reason: "staticfile",
		},
		State: files.State{
			Mode:          0644,
			Contents:      []byte(cfg.Content),
			FileExistence: files.FileExistencePresent,
		},
	}); err != nil {
		return nil, fmt.Errorf("failed to add file %q: %w", path, err)
	}
	return &ret, nil
}

func (f *Syncer) Name() config.Name {
	return "staticfile"
}

func (f *Syncer) Priority() drift.Priority {
	return drift.PriorityNormal
}

var _ drift.Detector = &Syncer{}
