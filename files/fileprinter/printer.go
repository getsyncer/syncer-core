package fileprinter

import (
	"io"

	"github.com/getsyncer/syncer-core/files"
)

type Printer interface {
	PrettyPrintDiffs(into io.Writer, toPrint *files.System[*files.DiffWithChangeReason]) error
}
