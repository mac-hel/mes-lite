package products_test

import (
	"errors"
	"testing"

	"github.com/mac-hel/mes-lite/internal/products"
)

func TestNewProduct(t *testing.T) {
	p := products.NewProduct(" VX-100 ", " Ventilation Unit X100 ", " piece ", products.CategoryVentilation)

	if p.SKU != "VX-100" {
		t.Errorf("SKU = %q, want %q", p.SKU, "VX-100")
	}
	if p.Name != "Ventilation Unit X100" {
		t.Errorf("Name = %q, want %q", p.Name, "Ventilation Unit X100")
	}
	if p.Unit != "piece" {
		t.Errorf("Unit = %q, want %q", p.Unit, "piece")
	}
	if p.Category != products.CategoryVentilation {
		t.Errorf("Category = %v, want %v", p.Category, products.CategoryVentilation)
	}
	if !p.IsActive {
		t.Error("IsActive should default to true")
	}
}

func TestProduct_UpdateDetails(t *testing.T) {
	p := products.NewProduct("VX-100", "Old Name", "piece", products.CategoryVentilation)

	p.UpdateDetails(" Updated Name ", " set ", products.CategoryFilter)

	if p.Name != "Updated Name" {
		t.Errorf("Name = %q, want %q", p.Name, "Updated Name")
	}
	if p.Unit != "set" {
		t.Errorf("Unit = %q, want %q", p.Unit, "set")
	}
	if p.Category != products.CategoryFilter {
		t.Errorf("Category = %v, want %v", p.Category, products.CategoryFilter)
	}
}

func TestProductCategory_String(t *testing.T) {
	tests := []struct {
		category products.ProductCategory
		want     string
	}{
		{products.CategoryVentilation, "Ventilation"},
		{products.CategoryFilter, "Filter"},
		{products.CategoryDuct, "Duct"},
		{products.CategoryMounting, "Mounting"},
		{products.CategoryOther, "Other"},
	}

	for _, tt := range tests {
		got := tt.category.String()
		if got != tt.want {
			t.Errorf("%d.String() = %q, want %q", tt.category, got, tt.want)
		}
	}
}

func TestProductCategory_String_Unknown(t *testing.T) {
	category := products.ProductCategory(99)

	if category.String() != "Unknown" {
		t.Errorf("String() = %q, want %q", category.String(), "Unknown")
	}
}

func TestProductCategory_Valid(t *testing.T) {
	tests := []struct {
		name     string
		category products.ProductCategory
		want     bool
	}{
		{"ventilation", products.CategoryVentilation, true},
		{"filter", products.CategoryFilter, true},
		{"duct", products.CategoryDuct, true},
		{"mounting", products.CategoryMounting, true},
		{"other", products.CategoryOther, true},
		{"negative", products.ProductCategory(-1), false},
		{"too high", products.ProductCategory(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product products.Product
		wantErr bool
	}{
		{"valid", products.NewProduct("VX-100", "Ventilation Unit", "piece", products.CategoryVentilation), false},
		{"blank sku", products.NewProduct(" ", "Ventilation Unit", "piece", products.CategoryVentilation), true},
		{"blank name", products.NewProduct("VX-100", " ", "piece", products.CategoryVentilation), true},
		{"blank unit", products.NewProduct("VX-100", "Ventilation Unit", " ", products.CategoryVentilation), true},
		{"invalid category", products.NewProduct("VX-100", "Ventilation Unit", "piece", products.ProductCategory(99)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if tt.wantErr {
				if !errors.Is(err, products.ErrInvalidProduct) {
					t.Fatalf("Validate() error = %v, want ErrInvalidProduct", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestProductCategory_ZeroValue(t *testing.T) {
	var c products.ProductCategory

	if c != products.CategoryVentilation {
		t.Errorf("zero value = %d, want %d (CategoryVentilation)", c, products.CategoryVentilation)
	}

	if c.String() != "Ventilation" {
		t.Errorf("zero value String() = %q, want %q", c.String(), "Ventilation")
	}
}

func TestProductCategory_IotaSequence(t *testing.T) {
	if products.CategoryVentilation != 0 {
		t.Errorf("CategoryVentilation = %d, want 0", products.CategoryVentilation)
	}
	if products.CategoryFilter != 1 {
		t.Errorf("CategoryFilter = %d, want 1", products.CategoryFilter)
	}
	if products.CategoryDuct != 2 {
		t.Errorf("CategoryDuct = %d, want 2", products.CategoryDuct)
	}
	if products.CategoryMounting != 3 {
		t.Errorf("CategoryMounting = %d, want 3", products.CategoryMounting)
	}
	if products.CategoryOther != 4 {
		t.Errorf("CategoryOther = %d, want 4", products.CategoryOther)
	}
}

func TestSentinelErrors(t *testing.T) {
	if products.ErrNotFound == nil {
		t.Error("ErrNotFound should not be nil")
	}
	if products.ErrNotFound.Error() != "product not found" {
		t.Errorf("ErrNotFound = %q, want %q", products.ErrNotFound.Error(), "product not found")
	}

	if products.ErrAlreadyExists == nil {
		t.Error("ErrAlreadyExists should not be nil")
	}
	if products.ErrAlreadyExists.Error() != "product already exists" {
		t.Errorf("ErrAlreadyExists = %q, want %q", products.ErrAlreadyExists.Error(), "product already exists")
	}

	if products.ErrInvalidProduct == nil {
		t.Error("ErrInvalidProduct should not be nil")
	}
	if products.ErrInvalidProduct.Error() != "invalid product" {
		t.Errorf("ErrInvalidProduct = %q, want %q", products.ErrInvalidProduct.Error(), "invalid product")
	}
}
