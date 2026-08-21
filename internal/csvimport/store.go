package csvimport

import (
	"context"

	"github.com/mac-hel/mes-lite/internal/production"
)

// Store persists validated CSV import records.
type Store interface {
	SaveBatch(ctx context.Context, records []ProductionEntryRecord) ([]production.Entry, error)
}
