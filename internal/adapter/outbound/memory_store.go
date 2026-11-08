package outbound

import (
	"context"
	"sync"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

// Compile-time check: MemoryStore must satisfy the port.
var _ port.LinkStore = (*MemoryStore)(nil)

type MemoryStore struct {
	mu    sync.RWMutex
	links map[string]domain.Link
}

func (m *MemoryStore) IncrementVisits(_ context.Context, code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.links[code]
	if !ok {
		return domain.ErrNotFound
	}
	l.Visits++
	m.links[code] = l
	return nil
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{links: make(map[string]domain.Link)}
}

func (m *MemoryStore) Save(_ context.Context, l domain.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.links[l.Code]; ok {
		return domain.ErrCodeTaken
	}
	m.links[l.Code] = l
	return nil
}

func (m *MemoryStore) FindByCode(_ context.Context, code string) (domain.Link, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.links[code]
	if !ok {
		return domain.Link{}, domain.ErrNotFound
	}
	return l, nil
}
