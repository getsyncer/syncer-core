package templatefiles

import (
	"context"
	"fmt"
	"text/template"

	"github.com/cresta/zapctx"
	"github.com/getsyncer/syncer-core/config"
	"github.com/getsyncer/syncer-core/drift"
	"github.com/getsyncer/syncer-core/drift/syncers/templatefiles/templatemutator"
	"github.com/getsyncer/syncer-core/files"
	"github.com/getsyncer/syncer-core/files/stateloader"
	"github.com/getsyncer/syncer-core/syncer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newGenerator[T config.TemplateConfig](files map[string]string, name config.Name, priority drift.Priority, decoder config.Decoder[T], logger *zapctx.Logger, setupLogic syncer.SetupSyncer, loader stateloader.StateLoader, postGenProcessor PostGenProcessor) (*Generator[T], error) {
	if name == "" {
		return nil, fmt.Errorf("name must be set")
	}
	generatedTemplates := make(map[string]*template.Template, len(files))
	for k, v := range files {
		tmpl, err := newTemplate(k, v)
		if err != nil {
			return nil, fmt.Errorf("unable to parse template %q: %w", k, err)
		}
		generatedTemplates[k] = tmpl
	}
	return &Generator[T]{
		files:            generatedTemplates,
		name:             name,
		priority:         priority,
		decoder:          decoder,
		logger:           logger,
		postGenProcessor: postGenProcessor,
		setupLogic:       setupLogic,
		loader:           loader,
	}, nil
}

type NewModuleConfig[T config.TemplateConfig] struct {
	Name             config.Name
	Files            map[string]string
	Priority         drift.Priority
	Decoder          config.Decoder[T]
	Setup            syncer.SetupSyncer
	PostGenProcessor PostGenProcessor
}

func NewModule[T config.TemplateConfig](newModuleConfig NewModuleConfig[T]) fx.Option {
	constructor := func(logger *zapctx.Logger, loader stateloader.StateLoader) (*Generator[T], error) {
		if newModuleConfig.Priority == 0 {
			newModuleConfig.Priority = drift.PriorityNormal
		}
		if newModuleConfig.Decoder == nil {
			newModuleConfig.Decoder = config.DefaultDecoder[T]()
		}
		if newModuleConfig.PostGenProcessor == nil {
			newModuleConfig.PostGenProcessor = PostGenProcessorList{}
		}
		return newGenerator(newModuleConfig.Files, newModuleConfig.Name, newModuleConfig.Priority, newModuleConfig.Decoder, logger, newModuleConfig.Setup, loader, newModuleConfig.PostGenProcessor)
	}
	return fx.Module(newModuleConfig.Name.String(),
		fx.Provide(
			fx.Annotate(
				constructor,
				fx.As(new(drift.Detector)),
				fx.ResultTags(drift.FxTagSyncers),
			),
		),
	)
}

type Generator[T config.TemplateConfig] struct {
	files            map[string]*template.Template
	name             config.Name
	priority         drift.Priority
	decoder          config.Decoder[T]
	mutators         templatemutator.MutatorList[T]
	setupLogic       syncer.SetupSyncer
	logger           *zapctx.Logger
	postGenProcessor PostGenProcessor
	loader           stateloader.StateLoader
}

var _ drift.Detector = &Generator[config.TemplateConfig]{}
var _ syncer.SetupSyncer = &Generator[config.TemplateConfig]{}

func (f *Generator[T]) Setup(ctx context.Context, runData *drift.RunData) error {
	if f.setupLogic != nil {
		return f.setupLogic.Setup(ctx, runData)
	}
	return nil
}

func (f *Generator[T]) AddMutator(mutator templatemutator.Mutator[T]) {
	f.mutators.AddMutator(mutator)
}

func (f *Generator[T]) PostGenProcess(ctx context.Context, fs *files.System[*files.StateWithChangeReason], runData *drift.RunData) error {
	return f.postGenProcessor.PostGenProcess(ctx, fs, runData)
}

func (f *Generator[T]) DetectDrift(ctx context.Context, runData *drift.RunData) (*files.System[*files.StateWithChangeReason], error) {
	f.logger.Debug(ctx, "running templatefile", zap.String("name", string(f.name)))
	cfg, err := f.decoder(runData.RunConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	cfg, err = f.mutators.Mutate(ctx, runData, f.loader, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to mutate config: %w", err)
	}
	var ret files.System[*files.StateWithChangeReason]
	for k, v := range f.files {
		f.logger.Debug(ctx, "generating template", zap.String("destination", k))
		var err error
		var fileContent string
		if fileContent, err = executeTemplateOnConfig(ctx, runData, cfg, v); err != nil {
			return nil, fmt.Errorf("unable to generate template for %s: %w", k, err)
		}
		if err := ret.Add(files.Path(k), &files.StateWithChangeReason{
			State: files.State{
				Contents:      []byte(fileContent),
				Mode:          0644,
				FileExistence: files.FileExistencePresent,
			},
			ChangeReason: &files.ChangeReason{
				Reason: "template",
			},
		}); err != nil {
			return nil, fmt.Errorf("unable to add file %s: %w", k, err)
		}
	}
	if err := f.PostGenProcess(ctx, &ret, runData); err != nil {
		return nil, fmt.Errorf("unable to post gen process: %w", err)
	}
	return &ret, nil
}

func (f *Generator[T]) Name() config.Name {
	return f.name
}

func (f *Generator[T]) Priority() drift.Priority {
	return f.priority
}
