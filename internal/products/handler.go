package products

import (
	"errors"
	"fmt"

	"github.com/go-fuego/fuego"
)

// Handler holds the HTTP handler methods for the products resource.
type Handler struct {
	store Store
}

// NewHandler creates a new Handler with the given Store.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// CreateProductRequest is the expected JSON body for creating a product.
// Validation tags enforce required fields.
type CreateProductRequest struct {
	SKU      string          `json:"sku"      validate:"required"`
	Name     string          `json:"name"     validate:"required"`
	Unit     string          `json:"unit"     validate:"required"`
	Category ProductCategory `json:"category" validate:"min=0,max=4"`
}

// Create handles POST /products and stores a new product.
func (h *Handler) Create(c fuego.ContextWithBody[CreateProductRequest]) (Product, error) {
	body, err := c.Body()
	if err != nil {
		return Product{}, err
	}

	p := NewProduct(body.SKU, body.Name, body.Unit, body.Category)
	if err := p.Validate(); err != nil {
		return Product{}, invalidProductError(err)
	}

	if err := h.store.Save(c.Context(), p); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return Product{}, fuego.ConflictError{
				Err:    err,
				Detail: fmt.Sprintf("product %q already exists", p.SKU),
			}
		}
		if errors.Is(err, ErrInvalidProduct) {
			return Product{}, invalidProductError(err)
		}
		return Product{}, err
	}

	return p, nil
}

// ListProductsResponse wraps a slice of products for JSON serialization.
type ListProductsResponse struct {
	Products []Product `json:"products"`
}

// List handles GET /products and returns all products.
func (h *Handler) List(c fuego.ContextNoBody) (ListProductsResponse, error) {
	prods, err := h.store.List(c.Context())
	if err != nil {
		return ListProductsResponse{}, err
	}

	if prods == nil {
		prods = []Product{}
	}

	return ListProductsResponse{Products: prods}, nil
}

// UpdateProductRequest is the expected JSON body for updating a product.
// Validation tags enforce required fields.
type UpdateProductRequest struct {
	Name     string          `json:"name"     validate:"required"`
	Unit     string          `json:"unit"     validate:"required"`
	Category ProductCategory `json:"category" validate:"min=0,max=4"`
}

// Update handles PUT /products/{sku} and replaces the product's mutable fields.
func (h *Handler) Update(c fuego.ContextWithBody[UpdateProductRequest]) (Product, error) {
	sku := c.PathParam("sku")

	p, err := h.store.FindBySKU(c.Context(), sku)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Product{}, fuego.NotFoundError{
				Err:    err,
				Detail: fmt.Sprintf("product %q not found", sku),
			}
		}
		return Product{}, err
	}

	body, err := c.Body()
	if err != nil {
		return Product{}, err
	}

	p.UpdateDetails(body.Name, body.Unit, body.Category)
	if err := p.Validate(); err != nil {
		return Product{}, invalidProductError(err)
	}

	if err := h.store.Update(c.Context(), p); err != nil {
		if errors.Is(err, ErrInvalidProduct) {
			return Product{}, invalidProductError(err)
		}
		return Product{}, err
	}

	return p, nil
}

func invalidProductError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{
		Err:    err,
		Detail: err.Error(),
	}
}

// Deactivate handles PUT /products/{sku}/deactivate and sets IsActive to false.
func (h *Handler) Deactivate(c fuego.ContextNoBody) (Product, error) {
	sku := c.PathParam("sku")

	p, err := h.store.FindBySKU(c.Context(), sku)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Product{}, fuego.NotFoundError{
				Err:    err,
				Detail: fmt.Sprintf("product %q not found", sku),
			}
		}
		return Product{}, err
	}

	p.IsActive = false

	if err := h.store.Update(c.Context(), p); err != nil {
		return Product{}, err
	}

	return p, nil
}

// Search handles GET /products/search?q=... and filters products by name or category.
func (h *Handler) Search(c fuego.ContextNoBody) (ListProductsResponse, error) {
	q := c.QueryParam("q")

	prods, err := h.store.Search(c.Context(), q)
	if err != nil {
		return ListProductsResponse{}, err
	}

	if prods == nil {
		prods = []Product{}
	}

	return ListProductsResponse{Products: prods}, nil
}
