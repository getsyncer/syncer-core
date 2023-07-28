package files

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath_Clean(t *testing.T) {
	require.Equal(t, Path("/foo/bar"), Path("/foo/bar").Clean())
	require.Equal(t, Path("/foo/bar"), Path("/foo//bar").Clean())
	require.Equal(t, Path("/foo/bar"), Path("/foo/./bar").Clean())
}
