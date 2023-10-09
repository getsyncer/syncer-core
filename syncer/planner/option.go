package planner

import "go.uber.org/fx"

type Option interface {
	apply(*plannerImpl)
}

type optionFunc func(*plannerImpl)

func (f optionFunc) apply(r *plannerImpl) {
	f(r)
}

func WithFilesAllowedNoMagicString(filenames ...string) Option {
	return optionFunc(func(r *plannerImpl) {
		r.filesAllowedNoMagicString = append(r.filesAllowedNoMagicString, filenames...)
	})
}

func FxOption(o Option) fx.Option {
	return fx.Provide(fx.Annotate(
		func() Option {
			return o
		},
		fx.ResultTags(`group:"planoption"`),
	))
	// Note: I tried both the below and neither worked.  I'm not sure why.
	//return fx.Supply(
	//	fx.Annotated{
	//		Group:  "planoption",
	//		Target: o,
	//	},
	//)
	//return fx.Supply(fx.Annotate(o, fx.ResultTags(`group:"planoption"`)))
}
