package applier

import (
	"context"
	"testing"

	"github.com/cresta/zapctx/testhelp/testhelp"
	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/diffexecutor/testdiffexecutor"
	"github.com/getsyncer/syncer-core/files/filestesthelp"
	"github.com/stretchr/testify/require"
)

func TestNewApplier(t *testing.T) {
	l := testhelp.ZapTestingLogger(t)
	de := testdiffexecutor.Test{}
	a := NewApplier(l, &de)
	ctx := context.Background()
	var f files.System[*files.DiffWithChangeReason]
	require.NoError(t, f.Add("a", &files.DiffWithChangeReason{
		Diff: filestesthelp.NewDiffNewFile("hello world"),
		ChangeReason: &files.ChangeReason{
			Reason: "test",
		},
	}))
	require.NoError(t, a.Apply(ctx, &f))
	require.Len(t, de.Executions, 1)
}
