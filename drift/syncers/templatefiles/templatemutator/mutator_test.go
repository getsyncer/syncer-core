package templatemutator

import (
	"context"
	"testing"

	"github.com/getsyncer/syncer-core/config"

	"github.com/stretchr/testify/require"

	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/files/stateloader"
)

type countingMutator[T config.TemplateConfig] struct {
	count int
	ret   T
}

func (c *countingMutator[T]) Mutate(_ context.Context, _ *drift.RunData, _ stateloader.StateLoader, _ T) (T, error) {
	c.count++
	return c.ret, nil
}

var _ Mutator[config.TemplateConfig] = &countingMutator[config.TemplateConfig]{}

type testConfig struct {
	name string
}

func TestMutatorList(t *testing.T) {
	cm := countingMutator[testConfig]{
		ret: testConfig{
			name: "replaced",
		},
	}
	ml := MutatorList[testConfig]{
		mutators: []Mutator[testConfig]{
			SimpleMutator[testConfig](func(in testConfig) (testConfig, error) {
				require.Equal(t, "original", in.name)
				return testConfig{
					name: "first_mutation",
				}, nil
			}),
			&cm,
		},
	}
	c, err := ml.Mutate(context.Background(), nil, nil, testConfig{name: "original"})
	require.NoError(t, err)
	require.Equal(t, "replaced", c.name)
	require.Equal(t, 1, cm.count)
}
