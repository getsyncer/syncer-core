package config

import "strings"

type Logic struct {
	Source string `yaml:"source"`
}

func (c *Logic) SourceWithoutVersion() string {
	parts := strings.SplitN(c.Source, "@", 2)
	if len(parts) == 1 {
		return c.Source
	}
	return parts[0]
}

func (c *Logic) SourceVersion() string {
	parts := strings.SplitN(c.Source, "@", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
