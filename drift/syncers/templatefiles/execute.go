package templatefiles

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"

	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/drift"
)

func executeTemplateOnConfig[T config.TemplateConfig](_ context.Context, runData *drift.RunData, config T, tmpl *template.Template) (string, error) {
	d := templateRunData[T]{
		RunData: runData,
		Config:  config,
	}
	var into bytes.Buffer
	if err := tmpl.Execute(&into, d); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return into.String(), nil
}

type templateRunData[T config.TemplateConfig] struct {
	RunData *drift.RunData
	Config  T
}

func newTemplate(name string, data string) (*template.Template, error) {
	tm := sprig.TxtFuncMap()
	tm["toYaml"] = toYAML
	return template.New(name).Funcs(tm).Parse(data)
}

// Taken from https://github.com/helm/helm/blob/main/pkg/engine/funcs.go#L30
func toYAML(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("unable to marshal to yaml: %w", err)
	}
	return strings.TrimSuffix(string(data), "\n"), nil
}
