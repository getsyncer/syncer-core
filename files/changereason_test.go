package files

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStateWithChangeReason(t *testing.T) {
	x := StateWithChangeReason{
		State: State{},
	}
	require.Error(t, x.Validate())
	x2 := StateWithChangeReason{
		State: State{
			Mode:          0644,
			Contents:      []byte("hello"),
			FileExistence: FileExistencePresent,
		},
	}
	require.NoError(t, x2.Validate())
}
