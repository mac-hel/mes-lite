package products

import "context"

// Store defines the persistence operations for Product data.
type Store interface {
	Save(ctx context.Context, p Product) error
	FindBySKU(ctx context.Context, sku string) (Product, error)
	List(ctx context.Context, opts ListOptions) ([]Product, error)
	Update(ctx context.Context, p Product) (Product, error)
}
