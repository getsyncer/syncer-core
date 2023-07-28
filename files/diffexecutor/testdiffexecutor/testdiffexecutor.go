package testdiffexecutor

import (
	"context"
	"fmt"

	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/diffexecutor"
)

type Test struct {
	Executions   []Executions
	FailingPaths []files.Path
}

func newTest() *Test {
	return &Test{}
}

type Executions struct {
	Path files.Path
	Diff *files.Diff
}

func (o *Test) ExecuteDiff(_ context.Context, path files.Path, d *files.Diff) error {
	for _, failingPath := range o.FailingPaths {
		if failingPath == path {
			return fmt.Errorf("failed to execute diff for %s", path)
		}
	}
	o.Executions = append(o.Executions, Executions{
		Path: path,
		Diff: d,
	})
	return nil
}

var _ diffexecutor.DiffExecutor = &Test{}
