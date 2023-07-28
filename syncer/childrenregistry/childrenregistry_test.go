package childrenregistry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChildrenRegistry(t *testing.T) {
	cc := ChildConfig{
		Content: []byte("hello"),
		Name:    "bob",
	}
	cr, err := New(cc)
	require.NoError(t, err)
	fetchedConfig, exists := cr.Get("bob")
	require.True(t, exists)
	require.Equal(t, cc, fetchedConfig)
	_, exists = cr.Get("alice")
	require.False(t, exists)
}
