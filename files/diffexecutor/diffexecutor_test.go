package diffexecutor_test

import (
	"context"
	"testing"

	"github.com/getsyncer/syncer-core/files/filestesthelp"

	"github.com/stretchr/testify/require"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/diffexecutor"
	"github.com/getsyncer/syncer-core/files/diffexecutor/testdiffexecutor"
)

func TestExecuteAllDiffs(t *testing.T) {
	var td testdiffexecutor.Test
	var fs files.System[*files.DiffWithChangeReason]
	diff := filestesthelp.NewDiffNewFile("hello world")
	require.NoError(t, fs.Add("a", &files.DiffWithChangeReason{
		Diff: diff,
	}))
	require.NoError(t, diffexecutor.ExecuteAllDiffs(context.Background(), &fs, &td))
	require.Equal(t, []testdiffexecutor.Executions{
		{
			Path: "a",
			Diff: diff,
		},
	}, td.Executions)
}
