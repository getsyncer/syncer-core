package fxregistry

import (
	"sync"

	"go.uber.org/fx"
)

type fxRegistry struct {
	options []fx.Option
	mu      sync.Mutex
}

var globalInstance fxRegistry

func (g *fxRegistry) Register(opt fx.Option) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.options = append(g.options, opt)
}

func (g *fxRegistry) Get() []fx.Option {
	g.mu.Lock()
	defer g.mu.Unlock()
	ret := make([]fx.Option, len(g.options))
	copy(ret, g.options)
	return ret
}

func Register(opt fx.Option) {
	globalInstance.Register(opt)
}

func Get() []fx.Option {
	return globalInstance.Get()
}
