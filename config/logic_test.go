package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogic_SourceVersion(t *testing.T) {
	testRun := func(source string, expectedVersion string) func(t *testing.T) {
		return func(t *testing.T) {
			logic := Logic{Source: source}
			require.Equal(t, expectedVersion, logic.SourceVersion())
		}
	}
	t.Run("no version", testRun("foo", ""))
	t.Run("with version", testRun("foo@bar", "bar"))
}

func TestLogic_SourceWithoutVersion(t *testing.T) {
	testRun := func(source string, expectedSource string) func(t *testing.T) {
		return func(t *testing.T) {
			logic := Logic{Source: source}
			require.Equal(t, expectedSource, logic.SourceWithoutVersion())
		}
	}
	t.Run("no version", testRun("foo", "foo"))
	t.Run("with version", testRun("foo@bar", "foo"))
}
