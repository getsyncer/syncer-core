package drifttest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/getsyncer/syncer-core/drift"

	"github.com/getsyncer/syncer-core/log/testlog"

	"go.uber.org/fx/fxtest"

	"github.com/getsyncer/syncer-core/fxcli"

	"github.com/getsyncer/syncer-core/syncerexec"

	"go.uber.org/fx"

	"github.com/getsyncer/syncer-core/osstub"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5"

	"github.com/getsyncer/syncer-core/files"
	"github.com/stretchr/testify/require"
)

type TestRunForConfig struct {
	fx.In
	Main     *syncerexec.RunSyncer
	TestStub osstub.OsStub
}

func ReasonableSampleFilesystem() *files.System[*files.State] {
	return files.NewSystem(map[files.Path]*files.State{
		"testfile_1": {
			Contents:      []byte("/dist\n/node_modules\n"),
			FileExistence: files.FileExistencePresent,
			Mode:          0644,
		},
		"testfile_2": {
			Contents:      []byte("package main\nfunc main() {}\n"),
			FileExistence: files.FileExistencePresent,
			Mode:          0644,
		},
	})
}

// TestOptions are only used by tests
func testOptions(t *testing.T) fx.Option {
	var allOpts []fx.Option
	allOpts = append(allOpts, osstub.TestModule, testlog.Module, fx.Supply(t), syncerexec.ModuleTest)
	return fx.Options(allOpts...)
}

func WithRun(config string, fs *files.System[*files.State], postRunVerification func(t *testing.T, items *Items)) func(t *testing.T) {
	return WithRunAndModule(config, fs, postRunVerification, fx.Module("empty"))
}

func WithRunAndModule(config string, fs *files.System[*files.State], postRunVerification func(t *testing.T, items *Items), module fx.Option) func(t *testing.T) {
	return func(t *testing.T) {
		trCon := TestRunConstructor(t, config, fs, postRunVerification)
		app := fxtest.New(t, module, syncerexec.DefaultFxOptions(), testOptions(t), fx.Provide(fx.Annotate(trCon, fx.As(new(fxcli.Main)))))
		app.RequireStart()
		<-app.Done()
	}
}

func TestRunConstructor(t *testing.T, config string, fs *files.System[*files.State], postRunVerification func(t *testing.T, items *Items)) func(cfg TestRunForConfig) *TestRun {
	return func(cfg TestRunForConfig) *TestRun {
		return &TestRun{
			T:                   t,
			Config:              config,
			Filesystem:          fs,
			TestRunForConfig:    &cfg,
			postRunVerification: postRunVerification,
		}
	}
}

type Items struct {
	TestRun *TestRun
}

type TestRun struct {
	TestRunForConfig    *TestRunForConfig
	T                   *testing.T
	Filesystem          *files.System[*files.State]
	Config              string
	tempDir             string
	postRunVerification func(t *testing.T, items *Items)
}

func (t *TestRun) SetupFiles() {
	require.Empty(t.T, t.tempDir)
	td := t.T.TempDir()
	t.tempDir = td
	require.NoError(t.T, os.Chdir(td))
	g, err := git.PlainInit(td, false)
	require.NoError(t.T, err)
	wt, err := g.Worktree()
	require.NoError(t.T, err)
	// Now make every file
	for _, path := range t.Filesystem.Paths() {
		fullPath := filepath.Clean(filepath.Join(td, string(path)))
		require.True(t.T, strings.HasPrefix(fullPath, td))
		filePath := t.Filesystem.Get(path)
		if filePath.FileExistence == files.FileExistenceAbsent {
			continue
		}
		m := filePath.Mode
		if m == 0 {
			m = 0644
		}
		require.NoError(t.T, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t.T, os.WriteFile(fullPath, t.Filesystem.Get(path).Contents, m))
		// now add to git
		_, err := wt.Add(string(path))
		require.NoError(t.T, err)
	}
	if t.Config != "" {
		require.NoError(t.T, os.MkdirAll(drift.DefaultSyncerGeneratedGoDirectory, 0755))
		configPath := filepath.Join(drift.DefaultSyncerGeneratedGoDirectory, drift.DefaultSyncerConfigFileName)
		require.NoError(t.T, os.WriteFile(configPath, []byte(t.Config), 0600))
		// now add to git
		_, err := wt.Add(configPath)
		require.NoError(t.T, err)
	}
	// Now commit
	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "noreply@example.com",
		},
	})
	require.NoError(t.T, err)
}

func (t *TestRun) Run() {
	t.SetupFiles()
	t.TestRunForConfig.TestStub.(*osstub.TestStub).Env["SYNCER_EXEC_CMD"] = "apply"
	t.TestRunForConfig.Main.Run()
	for _, p := range t.TestRunForConfig.TestStub.(*osstub.TestStub).Prints {
		t.T.Log(p)
	}
	if t.postRunVerification != nil {
		t.postRunVerification(t.T, &Items{TestRun: t})
	}
}

func (t *TestRun) MustExitCode(tes *testing.T, code int) {
	ts := t.TestRunForConfig.TestStub.(*osstub.TestStub)
	require.NotNil(tes, ts.ExitCode)
	require.Equal(tes, code, *ts.ExitCode)
}

func (t *TestRun) MustPrint(tes *testing.T, str string) {
	ts := t.TestRunForConfig.TestStub.(*osstub.TestStub)
	allPrints := strings.Join(ts.Prints, "\n")
	require.Contains(tes, allPrints, str)
}
