package reporting

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/mac-hel/mes-lite/internal/reporting/reportingdb"
)

// NewPostgresStore creates a PostgreSQL-backed reporting read store.
func NewPostgresStore(db reportingdb.DBTX) *PostgresStore {
	return &PostgresStore{queries: reportingdb.New(db)}
}

// PostgresStore reads reporting projections from PostgreSQL.
type PostgresStore struct {
	queries *reportingdb.Queries
}

// DailyProduction returns production quantities grouped by UTC day and product SKU.
func (s *PostgresStore) DailyProduction(ctx context.Context, from, to time.Time) ([]DailyProductionRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	rows, err := s.queries.DailyProduction(ctx, reportingdb.DailyProductionParams{
		FromTime: pgtype.Timestamptz{Time: from.UTC(), Valid: true},
		ToTime:   pgtype.Timestamptz{Time: to.UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	result := make([]DailyProductionRow, 0, len(rows))
	for _, row := range rows {
		if row.TotalQuantity > math.MaxInt || row.EntryCount > math.MaxInt {
			return nil, fmt.Errorf("report aggregate exceeds int size: %w", ErrInvalidRange)
		}
		result = append(result, DailyProductionRow{
			Day:           row.Day.Time.UTC(),
			ProductSKU:    row.ProductSku,
			TotalQuantity: int(row.TotalQuantity),
			EntryCount:    int(row.EntryCount),
		})
	}

	return result, nil
}

// EmployeeProductivity returns production quantities grouped by employee.
func (s *PostgresStore) EmployeeProductivity(ctx context.Context, from, to time.Time) ([]EmployeeProductivityRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	rows, err := s.queries.EmployeeProductivity(ctx, reportingdb.EmployeeProductivityParams{
		FromTime: pgtype.Timestamptz{Time: from.UTC(), Valid: true},
		ToTime:   pgtype.Timestamptz{Time: to.UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	result := make([]EmployeeProductivityRow, 0, len(rows))
	for _, row := range rows {
		if row.TotalQuantity > math.MaxInt || row.EntryCount > math.MaxInt {
			return nil, fmt.Errorf("report aggregate exceeds int size: %w", ErrInvalidRange)
		}
		result = append(result, EmployeeProductivityRow{
			EmployeeID:    row.EmployeeID,
			FirstName:     row.FirstName,
			LastName:      row.LastName,
			TotalQuantity: int(row.TotalQuantity),
			EntryCount:    int(row.EntryCount),
		})
	}

	return result, nil
}

// ProductStatistics returns production quantities grouped by product.
func (s *PostgresStore) ProductStatistics(ctx context.Context, from, to time.Time) ([]ProductStatisticsRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	rows, err := s.queries.ProductStatistics(ctx, reportingdb.ProductStatisticsParams{
		FromTime: pgtype.Timestamptz{Time: from.UTC(), Valid: true},
		ToTime:   pgtype.Timestamptz{Time: to.UTC(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	result := make([]ProductStatisticsRow, 0, len(rows))
	for _, row := range rows {
		if row.TotalQuantity > math.MaxInt || row.EntryCount > math.MaxInt || row.EmployeeCount > math.MaxInt {
			return nil, fmt.Errorf("report aggregate exceeds int size: %w", ErrInvalidRange)
		}
		result = append(result, ProductStatisticsRow{
			ProductSKU:    row.ProductSku,
			ProductName:   row.ProductName,
			TotalQuantity: int(row.TotalQuantity),
			EntryCount:    int(row.EntryCount),
			EmployeeCount: int(row.EmployeeCount),
		})
	}

	return result, nil
}
