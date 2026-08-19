package orders

import (
	"context"
	"fmt"
	"sync"
)

// NewInMemoryStore creates an in-memory [Store] for tests.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{orders: make(map[string]Order)}
}

// InMemoryStore stores production orders in memory.
type InMemoryStore struct {
	mu     sync.RWMutex
	orders map[string]Order
}

// Save stores a production order keyed by ID.
func (s *InMemoryStore) Save(ctx context.Context, order Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.orders[order.ID()]; exists {
		return fmt.Errorf("production order %q: %w", order.ID(), ErrAlreadyExists)
	}
	s.orders[order.ID()] = order
	return nil
}

// FindByID looks up a production order by ID.
func (s *InMemoryStore) FindByID(ctx context.Context, id string) (Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[id]
	if !ok {
		return Order{}, fmt.Errorf("production order %q: %w", id, ErrNotFound)
	}
	return order, nil
}

// Update replaces an existing production order and increments its version.
func (s *InMemoryStore) Update(ctx context.Context, order Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.orders[order.ID()]
	if !exists {
		return fmt.Errorf("production order %q: %w", order.ID(), ErrNotFound)
	}
	if current.Version() != order.Version() {
		return fmt.Errorf("production order %q version %d: %w", order.ID(), order.Version(), ErrVersionConflict)
	}
	s.orders[order.ID()] = order.incrementVersion()
	return nil
}
