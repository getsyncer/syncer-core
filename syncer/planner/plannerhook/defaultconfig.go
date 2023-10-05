package plannerhook

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"gopkg.in/yaml.v3"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/syncer/planner"
)

type DefaultConfig struct {
	ConfigToMerge config.Dynamic
}

func (d DefaultConfig) PreSetup(_ context.Context, rc *config.Root) error {
	if err := rc.Config.Merge(d.ConfigToMerge); err != nil {
		return fmt.Errorf("failed to merge default config: %w", err)
	}
	return nil
}

func NewDefaultConfig(configToMerge config.Dynamic) *DefaultConfig {
	return &DefaultConfig{
		ConfigToMerge: configToMerge,
	}
}

func NewDefaultConfigFromStruct(in interface{}) (*DefaultConfig, error) {
	yamlOut, err := yaml.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal yaml: %w", err)
	}
	var configToMerge config.Dynamic
	if err := yaml.Unmarshal(yamlOut, &configToMerge); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return NewDefaultConfig(configToMerge), nil
}

func DefaultConfigConstructor(in interface{}) func() (planner.Hook, error) {
	return func() (planner.Hook, error) {
		return NewDefaultConfigFromStruct(in)
	}
}

func DefaultConfigModule(name string, config interface{}) fx.Option {
	return planner.HookModule(name, DefaultConfigConstructor(config))
}

var _ planner.Hook = &DefaultConfig{}
