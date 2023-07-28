package files

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffResult_Validate(t *testing.T) {
	require.Error(t, (&DiffResult{
		DiffAction: DiffActionUnset,
	}).Validate())
	require.NoError(t, (&DiffResult{
		DiffAction: DiffActionNoChange,
	}).Validate())
}

func TestDiffAction_String(t *testing.T) {
	require.Equal(t, "unset", DiffActionUnset.String())
	require.Equal(t, "delete", DiffActionDelete.String())
	require.Equal(t, "create", DiffActionCreate.String())
	require.Equal(t, "update", DiffActionUpdate.String())
	require.Equal(t, "no change", DiffActionNoChange.String())
}
