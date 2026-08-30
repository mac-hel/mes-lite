package machines

import "context"

// Store persists received machine events.
type Store interface {
	Save(ctx context.Context, event Event) error
	FindByExternalEventID(ctx context.Context, machineID, externalEventID string) (Event, error)
	List(ctx context.Context) ([]Event, error)
}
