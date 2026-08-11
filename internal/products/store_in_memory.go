package products

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// NewInMemoryStore creates a map-based in-memory [InMemoryStore].
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{products: make(map[string]Product)}
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
	if _, ok := s.products[p.SKU]; ok {
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

// List returns products matching the given options.
func (s *InMemoryStore) List(_ context.Context, opts ListOptions) ([]Product, error) {
	opts, err := opts.normalize()
	if err != nil {
		return nil, err
	}
	prods := make([]Product, 0, len(s.products))
	for _, p := range s.products {
		if !matchesListOptions(p, opts) {
			continue
		}
		prods = append(prods, p)
	}
	sortProducts(prods, opts.Sort)
	return paginate(prods, opts), nil
}

// Update replaces the product at the given SKU and increments its version.
func (s *InMemoryStore) Update(_ context.Context, p Product) (Product, error) {
	if err := p.Validate(); err != nil {
		return Product{}, err
	}
	current, ok := s.products[p.SKU]
	if !ok {
		return Product{}, fmt.Errorf("product %q: %w", p.SKU, ErrNotFound)
	}
	if current.Version != p.Version {
		return Product{}, fmt.Errorf("product %q version %d: %w", p.SKU, p.Version, ErrVersionConflict)
	}
	p.Version++
	s.products[p.SKU] = p
	return p, nil
}

func matchesListOptions(p Product, opts ListOptions) bool {
	if opts.Active != nil && p.IsActive != *opts.Active {
		return false
	}
	if opts.Query == "" {
		return true
	}
	query := strings.ToLower(opts.Query)
	return strings.Contains(strings.ToLower(p.SKU), query) ||
		strings.Contains(strings.ToLower(p.Name), query) ||
		strings.Contains(strings.ToLower(p.Category.String()), query)
}

func sortProducts(prods []Product, sortKey string) {
	sort.Slice(prods, func(i, j int) bool {
		a, b := prods[i], prods[j]
		switch sortKey {
		case "-sku":
			return a.SKU > b.SKU
		case "name":
			return a.Name < b.Name
		case "-name":
			return a.Name > b.Name
		case "category":
			return a.Category < b.Category
		case "-category":
			return a.Category > b.Category
		default:
			return a.SKU < b.SKU
		}
	})
}

func paginate(prods []Product, opts ListOptions) []Product {
	if opts.Offset >= len(prods) {
		return []Product{}
	}
	end := opts.Offset + opts.Limit
	if end > len(prods) {
		end = len(prods)
	}
	return prods[opts.Offset:end]
}
