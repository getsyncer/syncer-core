package drift

import (
	"sort"
	"sync"

	"github.com/getsyncer/syncer-core/config"
)

type Registry interface {
	Registered() []Detector
	Get(name config.Name) (Detector, bool)
}

type registryImpl struct {
	syncers []Detector
	mu      sync.Mutex
}

func NewRegistry(syncers []Detector) (Registry, error) {
	seen := map[string]struct{}{}
	for _, s := range syncers {
		if _, ok := seen[string(s.Name())]; ok {
			return nil, &ErrSyncerAlreadyRegistered{Name: string(s.Name())}
		}
		seen[string(s.Name())] = struct{}{}
	}
	return &registryImpl{
		syncers: syncers,
	}, nil
}

func (r *registryImpl) Get(name config.Name) (Detector, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, s := range r.syncers {
		if s.Name() == name {
			return s, true
		}
	}
	return nil, false
}

var _ Registry = &registryImpl{}

type ErrSyncerAlreadyRegistered struct {
	Name string
}

func (e *ErrSyncerAlreadyRegistered) Error() string {
	return "syncer already registered: " + e.Name
}

func (r *registryImpl) Registered() []Detector {
	r.mu.Lock()
	defer r.mu.Unlock()
	sort.SliceStable(r.syncers, func(i, j int) bool {
		return r.syncers[i].Priority() < r.syncers[j].Priority() || r.syncers[i].Name() < r.syncers[j].Name()
	})
	return r.syncers
}
