package childrenregistry

import (
	"github.com/getsyncer/syncer-core/config"
	"go.uber.org/fx"
)

const (
	FxTagChildren = `group:"childrensource"`
)

func NewModule(name config.Name, content []byte) fx.Option {
	constructor := func() ChildConfig {
		return ChildConfig{
			Content: content,
			Name:    name,
		}
	}
	return fx.Module(name.String(),
		fx.Provide(
			fx.Annotate(
				constructor,
				fx.ResultTags(FxTagChildren),
			),
		),
	)
}

var Module = fx.Module("childregistry",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(ChildrenRegistry)),
			fx.ParamTags(FxTagChildren),
		),
	),
)
