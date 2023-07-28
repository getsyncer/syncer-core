package drifttest_test

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/getsyncer/syncer-core/drifttest"
)

func TestTestRun(t *testing.T) {
	x := drifttest.TestRun{
		T:          t,
		Filesystem: drifttest.ReasonableSampleFilesystem(),
	}
	x.SetupFiles()
	contents, err := os.ReadFile("testfile_1")
	require.NoError(t, err)
	require.Equal(t, "/dist\n/node_modules\n", string(contents))
	cmd := exec.Command("git", "ls-files")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	require.NoError(t, cmd.Run())
	require.Equal(t, "testfile_1\ntestfile_2\n", stdout.String())
}

func TestNoSyncFile(t *testing.T) {
	t.Run("missingconfig", drifttest.WithRun("", drifttest.ReasonableSampleFilesystem(), func(t *testing.T, items *drifttest.Items) {
		items.TestRun.MustExitCode(t, 1)
		items.TestRun.MustPrint(t, "no config file found")
	}))
}

func TestBadVersion(t *testing.T) {
	t.Run("bad-config", drifttest.WithRun("version: 3", drifttest.ReasonableSampleFilesystem(), func(t *testing.T, items *drifttest.Items) {
		items.TestRun.MustExitCode(t, 1)
		items.TestRun.MustPrint(t, "unknown config version")
	}))
}

func TestEmptySync(t *testing.T) {
	t.Run("empty-run", drifttest.WithRun("version: 1", drifttest.ReasonableSampleFilesystem(), func(t *testing.T, items *drifttest.Items) {
		items.TestRun.MustExitCode(t, 0)
	}))
}
