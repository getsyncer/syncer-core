package syncer

import (
	"context"
	"testing"

	"github.com/getsyncer/syncer-core/drift"

	"github.com/stretchr/testify/require"
)

func TestSetupSyncerFunc_Setup(t *testing.T) {
	i := 0
	s1 := SetupSyncerFunc(func(_ context.Context, _ *drift.RunData) error {
		i++
		return nil
	})
	j := 0
	s2 := SetupSyncerFunc(func(_ context.Context, _ *drift.RunData) error {
		j++
		return nil
	})
	m := MultiSetupSyncer([]SetupSyncer{s1, s2})
	require.NoError(t, m.Setup(context.Background(), nil))
	require.Equal(t, 1, i)
	require.Equal(t, 1, j)
}
