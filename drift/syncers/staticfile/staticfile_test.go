package staticfile

import (
	"testing"

	"github.com/getsyncer/syncer-core/drifttest"
	"github.com/getsyncer/syncer-core/files"
)

func TestStaticFileDrift(t *testing.T) {
	config := `
version: 1
logic:
  - source: github.com/getsyncer/syncer-core/drift/syncers/staticfile
syncs:
  - logic: staticfile
    config:
      filename: testfile_1
      content: Static sync content
`
	t.Run("update-existing-file", drifttest.WithRun(config, drifttest.ReasonableSampleFilesystem(), func(t *testing.T, items *drifttest.Items) {
		items.TestRun.MustExitCode(t, 0)
		drifttest.FileContains(t, "testfile_1", "Static sync content")
		drifttest.OnlyGitChanges(t, "testfile_1")
	}))
	t.Run("no-change", drifttest.WithRun(config, files.SimpleState(map[string]string{
		"testfile_1": "Static sync content",
	}), func(t *testing.T, items *drifttest.Items) {
		items.TestRun.MustExitCode(t, 0)
		drifttest.FileContains(t, "testfile_1", "Static sync content")
		drifttest.OnlyGitChanges(t)
	}))
	t.Run("make-new-file", drifttest.WithRun(config, files.SimpleState(map[string]string{
		"testfile_2": "Static sync content",
	}), func(t *testing.T, items *drifttest.Items) {
		items.TestRun.MustExitCode(t, 0)
		drifttest.FileContains(t, "testfile_1", "Static sync content")
		drifttest.OnlyGitChanges(t, "testfile_1")
	}))
}
