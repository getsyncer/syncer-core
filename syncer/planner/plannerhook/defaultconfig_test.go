package plannerhook_test

import (
	"testing"

	"github.com/getsyncer/syncer-core/syncer/planner/plannerhook"
	"go.uber.org/fx"

	_ "github.com/getsyncer/syncer-core/drift/syncers/staticfile"
	"github.com/getsyncer/syncer-core/drifttest"
)

func TestDefaultConfig(t *testing.T) {
	config := `
version: 1
logic:
  - source: github.com/getsyncer/syncer-core/drift/syncers/staticfile
syncs:
  - logic: staticfile
    config:
      filename: testfile_1
`
	type exampleConfig struct {
		Content string `yaml:"content"`
	}
	testModule := fx.Module("test", plannerhook.DefaultConfigModule("tests", exampleConfig{Content: "Static sync content"}))
	t.Run("update-existing-file", drifttest.WithRunAndModule(config, drifttest.ReasonableSampleFilesystem(), func(t *testing.T, items *drifttest.Items) {
		items.TestRun.MustExitCode(t, 0)
		drifttest.FileContains(t, "testfile_1", "Static sync content")
		drifttest.OnlyGitChanges(t, "testfile_1")
	}, testModule))
}
