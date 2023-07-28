package fxregistry

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestRegistry(t *testing.T) {
	var f fxRegistry
	require.Len(t, f.Get(), 0)
	m := fx.Module("testing")
	f.Register(m)
	require.Len(t, f.Get(), 1)
	require.Equal(t, m, f.Get()[0])
}
