package reporting

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrInvalidRange is returned when a report time range cannot produce a meaningful result.
var ErrInvalidRange = errors.New("invalid report range")

// Store reads reporting data. It intentionally exposes read models, not domain aggregates.
type Store interface {
	DailyProduction(ctx context.Context, from, to time.Time) ([]DailyProductionRow, error)
}

// DailyProductionRow is a read model for production totals grouped by UTC day and product.
type DailyProductionRow struct {
	Day           time.Time
	ProductSKU    string
	TotalQuantity int
	EntryCount    int
}

func validateRange(from, to time.Time) error {
	if from.IsZero() {
		return fmt.Errorf("from is required: %w", ErrInvalidRange)
	}
	if to.IsZero() {
		return fmt.Errorf("to is required: %w", ErrInvalidRange)
	}
	if !from.Before(to) {
		return fmt.Errorf("from must be before to: %w", ErrInvalidRange)
	}

	return nil
}
