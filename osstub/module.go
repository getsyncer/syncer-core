package osstub

import "go.uber.org/fx"

var Module = fx.Module("osstub",
	fx.Supply(fx.Annotate(&osStub{}, fx.As(new(OsStub)))),
)

var TestModule = fx.Module("osstub-test",
	fx.Supply(fx.Annotate(&TestStub{
		Env: make(map[string]string),
	}, fx.As(new(OsStub)))),
)
