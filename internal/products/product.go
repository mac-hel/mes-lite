package products

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNotFound is returned when a product cannot be found by SKU.
var ErrNotFound = errors.New("product not found")

// ErrAlreadyExists is returned when trying to create a product with a duplicate SKU.
var ErrAlreadyExists = errors.New("product already exists")

// ErrInvalidProduct is returned when product data breaks domain rules.
var ErrInvalidProduct = errors.New("invalid product")

// ErrVersionConflict is returned when an update uses a stale product version.
var ErrVersionConflict = errors.New("product version conflict")

// ProductCategory classifies products by their function.
type ProductCategory int

// Product categories.
const (
	CategoryVentilation ProductCategory = iota
	CategoryFilter
	CategoryDuct
	CategoryMounting
	CategoryOther
)

func (c ProductCategory) String() string {
	switch c {
	case CategoryVentilation:
		return "Ventilation"
	case CategoryFilter:
		return "Filter"
	case CategoryDuct:
		return "Duct"
	case CategoryMounting:
		return "Mounting"
	case CategoryOther:
		return "Other"
	default:
		return "Unknown"
	}
}

// Valid reports whether the category is one of the supported product categories.
func (c ProductCategory) Valid() bool {
	return c >= CategoryVentilation && c <= CategoryOther
}

// Product represents a manufactured ventilation component.
type Product struct {
	SKU      string
	Name     string
	Category ProductCategory
	Unit     string
	IsActive bool
	Version  int
}

// NewProduct creates a valid Product with IsActive set to true.
func NewProduct(sku, name, unit string, category ProductCategory) (Product, error) {
	p := Product{
		SKU:      strings.TrimSpace(sku),
		Name:     strings.TrimSpace(name),
		Unit:     strings.TrimSpace(unit),
		Category: category,
		IsActive: true,
		Version:  1,
	}
	if err := p.Validate(); err != nil {
		return Product{}, err
	}

	return p, nil
}

// UpdateDetails replaces mutable product fields and preserves product invariants.
func (p *Product) UpdateDetails(name, unit string, category ProductCategory) error {
	updated := *p
	updated.Name = strings.TrimSpace(name)
	updated.Unit = strings.TrimSpace(unit)
	updated.Category = category
	if err := updated.Validate(); err != nil {
		return err
	}

	*p = updated
	return nil
}

// Validate checks the product invariants that must hold in every entry point.
func (p Product) Validate() error {
	if strings.TrimSpace(p.SKU) == "" {
		return fmt.Errorf("sku is required: %w", ErrInvalidProduct)
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("name is required: %w", ErrInvalidProduct)
	}
	if strings.TrimSpace(p.Unit) == "" {
		return fmt.Errorf("unit is required: %w", ErrInvalidProduct)
	}
	if !p.Category.Valid() {
		return fmt.Errorf("category %d is not supported: %w", p.Category, ErrInvalidProduct)
	}
	if p.Version <= 0 {
		return fmt.Errorf("version must be greater than zero: %w", ErrInvalidProduct)
	}

	return nil
}
