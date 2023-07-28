package drift_test

import (
	"testing"

	"github.com/getsyncer/syncer-core/drift"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/syncer/syncertesthelper"
	"github.com/stretchr/testify/require"
)

func testRegistry(t *testing.T, r drift.Registry, mustExist config.Name) {
	_, exists := r.Get("unknown")
	require.False(t, exists)
	ds, exists := r.Get(mustExist)
	require.True(t, exists)
	require.Equal(t, mustExist, ds.Name())
}

func TestRegistry(t *testing.T) {
	testSyncer := syncertesthelper.TestSyncer{
		Ret:        nil,
		SyncerName: "test1",
	}
	testSyncer2 := syncertesthelper.TestSyncer{
		Ret:        nil,
		SyncerName: "test2",
	}
	reg, err := drift.NewRegistry([]drift.Detector{&testSyncer, &testSyncer2})
	require.NoError(t, err)
	require.Len(t, reg.Registered(), 2)
	testRegistry(t, reg, "test1")
	testRegistry(t, reg, "test2")
}
