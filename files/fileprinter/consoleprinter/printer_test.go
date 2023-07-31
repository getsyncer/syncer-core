package consoleprinter

import (
	"bytes"
	"testing"

	"github.com/getsyncer/syncer-core/files/filestesthelp"

	"github.com/getsyncer/syncer-core/files"
	"github.com/stretchr/testify/require"
)

func TestConsolePrinterImpl_PrettyPrintDiffs(t *testing.T) {
	var x consolePrinterImpl
	var into bytes.Buffer
	var toPrint files.System[*files.DiffWithChangeReason]
	require.NoError(t, x.PrettyPrintDiffs(&into, &toPrint))
	require.NoError(t, toPrint.Add("foo", &files.DiffWithChangeReason{
		Diff: filestesthelp.NewDiffDeleteFile(nil),
		ChangeReason: &files.ChangeReason{
			Reason: "foo",
		},
	}))
	require.NoError(t, x.PrettyPrintDiffs(&into, &toPrint))
	require.Contains(t, into.String(), "foo")
}
