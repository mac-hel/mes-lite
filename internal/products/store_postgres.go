package products

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mac-hel/mes-lite/internal/products/productsdb"
)

// NewPostgresStore creates a PostgreSQL-backed [Store].
func NewPostgresStore(db productsdb.DBTX) *PostgresStore {
	return &PostgresStore{queries: productsdb.New(db)}
}

// PostgresStore stores products in PostgreSQL through sqlc-generated queries.
type PostgresStore struct {
	queries *productsdb.Queries
}

// Save stores a product keyed by SKU.
func (s *PostgresStore) Save(ctx context.Context, p Product) error {
	if err := p.Validate(); err != nil {
		return err
	}
	_, err := s.queries.CreateProduct(ctx, productsdb.CreateProductParams{
		Sku:      p.SKU,
		Name:     p.Name,
		Category: int32(p.Category),
		Unit:     p.Unit,
		IsActive: p.IsActive,
		Version:  int32(p.Version),
	})
	if err != nil {
		return mapPostgresError(p.SKU, err)
	}
	return nil
}

// FindBySKU looks up a product by SKU. Returns [ErrNotFound] if not found.
func (s *PostgresStore) FindBySKU(ctx context.Context, sku string) (Product, error) {
	p, err := s.queries.GetProduct(ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, fmt.Errorf("product %q: %w", sku, ErrNotFound)
		}
		return Product{}, err
	}
	return productFromGetDB(p), nil
}

// List returns products matching the given options.
func (s *PostgresStore) List(ctx context.Context, opts ListOptions) ([]Product, error) {
	opts, err := opts.normalize()
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.ListProducts(ctx, productsdb.ListProductsParams{
		Query:       opts.Query,
		Active:      activeFilter(opts.Active),
		Sort:        opts.Sort,
		OffsetValue: int32(opts.Offset),
		LimitValue:  int32(opts.Limit),
	})
	if err != nil {
		return nil, err
	}
	prods := make([]Product, 0, len(rows))
	for _, row := range rows {
		prods = append(prods, productFromListDB(row))
	}
	return prods, nil
}

// Update replaces the product at the given SKU and increments its version.
func (s *PostgresStore) Update(ctx context.Context, p Product) (Product, error) {
	if err := p.Validate(); err != nil {
		return Product{}, err
	}
	row, err := s.queries.UpdateProduct(ctx, productsdb.UpdateProductParams{
		Sku:      p.SKU,
		Name:     p.Name,
		Category: int32(p.Category),
		Unit:     p.Unit,
		IsActive: p.IsActive,
		Version:  int32(p.Version),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return s.updateNoRowsError(ctx, p.SKU, p.Version)
		}
		return Product{}, mapPostgresError(p.SKU, err)
	}
	return productFromUpdateDB(row), nil
}

func productFromGetDB(p productsdb.GetProductRow) Product {
	return Product{SKU: p.Sku, Name: p.Name, Category: ProductCategory(p.Category), Unit: p.Unit, IsActive: p.IsActive, Version: int(p.Version)}
}

func productFromListDB(p productsdb.ListProductsRow) Product {
	return Product{SKU: p.Sku, Name: p.Name, Category: ProductCategory(p.Category), Unit: p.Unit, IsActive: p.IsActive, Version: int(p.Version)}
}

func productFromUpdateDB(p productsdb.UpdateProductRow) Product {
	return Product{SKU: p.Sku, Name: p.Name, Category: ProductCategory(p.Category), Unit: p.Unit, IsActive: p.IsActive, Version: int(p.Version)}
}

func (s *PostgresStore) updateNoRowsError(ctx context.Context, sku string, version int) (Product, error) {
	_, err := s.FindBySKU(ctx, sku)
	if err != nil {
		return Product{}, err
	}
	return Product{}, fmt.Errorf("product %q version %d: %w", sku, version, ErrVersionConflict)
}

func mapPostgresError(sku string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("product %q: %w", sku, ErrAlreadyExists)
		case "23514", "23502":
			return fmt.Errorf("product %q: %w", sku, ErrInvalidProduct)
		}
	}
	return err
}
