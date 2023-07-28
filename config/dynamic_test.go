package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRun_Decode(t *testing.T) {
	data := `
name: test
user:
  age: 10
`
	type User struct {
		Age int    `yaml:"age,omitempty"`
		Job string `yaml:"job,omitempty"`
	}
	type Config struct {
		Name string `yaml:"name,omitempty"`
		User User   `yaml:"user,omitempty"`
	}
	c := &Config{
		Name: "test",
		User: User{
			Age: 10,
		},
	}
	var r Dynamic
	require.NoError(t, yaml.Unmarshal([]byte(data), &r))
	var c2 Config
	require.NoError(t, r.Decode(&c2))
	require.Equal(t, c, &c2)
	var r2 Dynamic
	c3 := &Config{
		User: User{
			Job: "testjob",
		},
	}
	c3Data, err := yaml.Marshal(c3)
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(c3Data, &r2))
	require.NoError(t, r2.Merge(r))
	var c4 Config
	require.NoError(t, r2.Decode(&c4))
	require.Equal(t, &Config{
		Name: "test",
		User: User{
			Age: 10,
			Job: "testjob",
		},
	}, &c4)
	r2Yaml, err := r2.AsYaml()
	require.NoError(t, err)
	var c5 Config
	require.NoError(t, yaml.Unmarshal([]byte(r2Yaml), &c5))
	require.Equal(t, c4, c5)
}
