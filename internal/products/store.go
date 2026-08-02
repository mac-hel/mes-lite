package products

import (
	"context"
	"fmt"
	"strings"
)

// Store defines the persistence operations for Product data.
type Store interface {
	Save(ctx context.Context, p Product) error
	FindBySKU(ctx context.Context, sku string) (Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, p Product) error
	Search(ctx context.Context, query string) ([]Product, error)
}

// NewInMemoryStore creates a map-based in-memory [InMemoryStore].
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		products: make(map[string]Product),
	}
}

// InMemoryStore is a map-based in-memory implementation of [Store].
type InMemoryStore struct {
	products map[string]Product
}

// Save stores a product keyed by SKU.
func (s *InMemoryStore) Save(_ context.Context, p Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	_, ok := s.products[p.SKU]
	if ok {
		return fmt.Errorf("product %q: %w", p.SKU, ErrAlreadyExists)
	}
	s.products[p.SKU] = p
	return nil
}

// FindBySKU looks up a product by SKU. Returns [ErrNotFound] if not found.
func (s *InMemoryStore) FindBySKU(_ context.Context, sku string) (Product, error) {
	p, ok := s.products[sku]
	if !ok {
		return Product{}, fmt.Errorf("product %q: %w", sku, ErrNotFound)
	}
	return p, nil
}

// List returns all products.
func (s *InMemoryStore) List(_ context.Context) ([]Product, error) {
	prods := make([]Product, 0, len(s.products))
	for _, p := range s.products {
		prods = append(prods, p)
	}
	return prods, nil
}

// Update replaces the product at the given SKU. Returns [ErrNotFound] if not found.
func (s *InMemoryStore) Update(_ context.Context, p Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	if _, ok := s.products[p.SKU]; !ok {
		return fmt.Errorf("product %q: %w", p.SKU, ErrNotFound)
	}
	s.products[p.SKU] = p
	return nil
}

// Search filters products whose name or category contains the query (case-insensitive).
func (s *InMemoryStore) Search(_ context.Context, query string) ([]Product, error) {
	query = strings.ToLower(query)
	prods := make([]Product, 0)
	for _, p := range s.products {
		if strings.Contains(strings.ToLower(p.Name), query) || strings.Contains(strings.ToLower(p.Category.String()), query) {
			prods = append(prods, p)
		}
	}
	return prods, nil
}
