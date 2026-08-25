package production

import (
	"context"
)

// EntryStore defines persistence operations for production entry data.
type EntryStore interface {
	Save(ctx context.Context, entry Entry) error
	FindByID(ctx context.Context, id string) (Entry, error)
	FindByRequestID(ctx context.Context, requestID string) (Entry, error)
	List(ctx context.Context, opts ListOptions) ([]Entry, error)
}

// CorrectionStore defines persistence operations for production-entry corrections.
type CorrectionStore interface {
	SaveCorrection(ctx context.Context, correction Correction) error
	ListCorrections(ctx context.Context, entryID string) ([]Correction, error)
}
