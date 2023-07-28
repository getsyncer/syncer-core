package fxcli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/fx"

	"go.uber.org/fx/fxtest"
)

type testApp struct {
	runCount int
}

func (t *testApp) Run() {
	t.runCount++
}

var _ Main = (*testApp)(nil)

func TestNewFxCLI(t *testing.T) {
	inst := &testApp{}
	app := fxtest.New(t, Module, fx.Supply(fx.Annotate(inst, fx.As(new(Main)))))
	app.RequireStart()
	app.RequireStop()
	require.Equal(t, 1, inst.runCount)
}
