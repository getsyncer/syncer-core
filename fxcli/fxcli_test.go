package fxcli

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/fx"

	"go.uber.org/fx/fxtest"
)

type testApp struct {
	runCount atomic.Int32
}

func (t *testApp) Run() {
	t.runCount.Add(1)
}

var _ Main = (*testApp)(nil)

func TestNewFxCLI(t *testing.T) {
	inst := &testApp{}
	app := fxtest.New(t, Module, fx.Supply(fx.Annotate(inst, fx.As(new(Main)))))
	app.RequireStart()
	app.RequireStop()
	require.Equal(t, int32(1), inst.runCount.Load())
}
