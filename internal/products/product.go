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
}

// NewProduct creates a new Product with IsActive set to true.
func NewProduct(sku, name, unit string, category ProductCategory) Product {
	return Product{
		SKU:      strings.TrimSpace(sku),
		Name:     strings.TrimSpace(name),
		Unit:     strings.TrimSpace(unit),
		Category: category,
		IsActive: true,
	}
}

// UpdateDetails replaces mutable product fields and applies the same normalization as creation.
func (p *Product) UpdateDetails(name, unit string, category ProductCategory) {
	p.Name = strings.TrimSpace(name)
	p.Unit = strings.TrimSpace(unit)
	p.Category = category
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

	return nil
}
