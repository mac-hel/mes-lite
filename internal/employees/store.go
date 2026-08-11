package employees

import (
	"context"
)

// Store defines the persistence operations for Employee data.
type Store interface {
	Save(ctx context.Context, emp Employee) error
	FindByID(ctx context.Context, id string) (Employee, error)
	List(ctx context.Context, opts ListOptions) ([]Employee, error)
	Update(ctx context.Context, emp Employee) (Employee, error)
}
