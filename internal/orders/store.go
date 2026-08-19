package orders

import "context"

// Store defines persistence operations for production orders.
type Store interface {
	Save(ctx context.Context, order Order) error
	FindByID(ctx context.Context, id string) (Order, error)
	Update(ctx context.Context, order Order) error
}
