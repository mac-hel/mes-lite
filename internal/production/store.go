package production

import (
	"context"
)

// Store defines the persistence operations for production entry data.
type Store interface {
	Save(ctx context.Context, entry Entry) error
	FindByID(ctx context.Context, id string) (Entry, error)
	List(ctx context.Context, opts ListOptions) ([]Entry, error)
}
