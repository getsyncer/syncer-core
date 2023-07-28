package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoot_AsYaml(t *testing.T) {
	testRun := func(r *Root, expectedYaml string) func(t *testing.T) {
		return func(t *testing.T) {
			yaml, err := r.AsYaml()
			require.NoError(t, err)
			require.Equal(t, expectedYaml, yaml)
		}
	}
	t.Run("empty", testRun(&Root{}, "version: 0\n"))
}

func TestRoot_Merge(t *testing.T) {
	testRun := func(r1, r2 *Root, expected *Root) func(t *testing.T) {
		return func(t *testing.T) {
			err := r1.Merge(r2)
			require.NoError(t, err)
			require.Equal(t, expected, r1)
		}
	}
	t.Run("empty", testRun(&Root{}, &Root{}, &Root{}))
	t.Run("general", testRun(
		&Root{
			Version: 1,
			Children: []Logic{
				{
					Source: "a",
				},
				{
					Source: "b",
				},
			},
		},
		&Root{
			Version: 1,
			Children: []Logic{
				{
					Source: "a",
				},
				{
					Source: "c",
				},
			},
		},
		&Root{
			Version: 1,
			Children: []Logic{
				{Source: "a"},
				{Source: "b"},
				{Source: "c"},
			},
		}))
}
