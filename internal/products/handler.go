package products

import (
	"errors"
	"fmt"
	"strconv"

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

	p, err := NewProduct(body.SKU, body.Name, body.Unit, body.Category)
	if err != nil {
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
	Products   []Product `json:"products"`
	Pagination Page      `json:"pagination"`
}

// List handles GET /products and returns all products.
func (h *Handler) List(c fuego.ContextNoBody) (ListProductsResponse, error) {
	opts, err := listOptionsFromQuery(c)
	if err != nil {
		return ListProductsResponse{}, invalidListOptionsError(err)
	}

	prods, err := h.store.List(c.Context(), opts)
	if err != nil {
		if errors.Is(err, ErrInvalidListOptions) {
			return ListProductsResponse{}, invalidListOptionsError(err)
		}
		return ListProductsResponse{}, err
	}

	if prods == nil {
		prods = []Product{}
	}

	return ListProductsResponse{Products: prods, Pagination: Page{Limit: opts.Limit, Offset: opts.Offset, Count: len(prods)}}, nil
}

func listOptionsFromQuery(c fuego.ContextNoBody) (ListOptions, error) {
	limit, err := parseIntQuery(c.QueryParam("limit"), "limit")
	if err != nil {
		return ListOptions{}, err
	}
	offset, err := parseIntQuery(c.QueryParam("offset"), "offset")
	if err != nil {
		return ListOptions{}, err
	}
	active, err := parseBoolQuery(c.QueryParam("active"), "active")
	if err != nil {
		return ListOptions{}, err
	}

	return ListOptions{Limit: limit, Offset: offset, Sort: c.QueryParam("sort"), Query: c.QueryParam("q"), Active: active}.normalize()
}

func parseIntQuery(raw, name string) (int, error) {
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, ErrInvalidListOptions)
	}
	return value, nil
}

func parseBoolQuery(raw, name string) (*bool, error) {
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fmt.Errorf("%s must be true or false: %w", name, ErrInvalidListOptions)
	}
	return &value, nil
}

func invalidListOptionsError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{Err: err, Detail: err.Error()}
}

// UpdateProductRequest is the expected JSON body for updating a product.
// Validation tags enforce required fields.
type UpdateProductRequest struct {
	Name     string          `json:"name"     validate:"required"`
	Unit     string          `json:"unit"     validate:"required"`
	Category ProductCategory `json:"category" validate:"min=0,max=4"`
	Version  int             `json:"version"  validate:"min=1"`
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

	if err := p.UpdateDetails(body.Name, body.Unit, body.Category); err != nil {
		return Product{}, invalidProductError(err)
	}
	p.Version = body.Version

	updated, err := h.store.Update(c.Context(), p)
	if err != nil {
		if errors.Is(err, ErrInvalidProduct) {
			return Product{}, invalidProductError(err)
		}
		if errors.Is(err, ErrVersionConflict) {
			return Product{}, versionConflictError(err)
		}
		return Product{}, err
	}

	return updated, nil
}

func invalidProductError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{
		Err:    err,
		Detail: err.Error(),
	}
}

func versionConflictError(err error) fuego.ConflictError {
	return fuego.ConflictError{Err: err, Detail: err.Error()}
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

	updated, err := h.store.Update(c.Context(), p)
	if err != nil {
		if errors.Is(err, ErrVersionConflict) {
			return Product{}, versionConflictError(err)
		}
		return Product{}, err
	}

	return updated, nil
}

// Search handles GET /products/search?q=... and filters products by name or category.
func (h *Handler) Search(c fuego.ContextNoBody) (ListProductsResponse, error) {
	opts, err := listOptionsFromQuery(c)
	if err != nil {
		return ListProductsResponse{}, invalidListOptionsError(err)
	}

	prods, err := h.store.List(c.Context(), opts)
	if err != nil {
		if errors.Is(err, ErrInvalidListOptions) {
			return ListProductsResponse{}, invalidListOptionsError(err)
		}
		return ListProductsResponse{}, err
	}

	if prods == nil {
		prods = []Product{}
	}

	return ListProductsResponse{Products: prods, Pagination: Page{Limit: opts.Limit, Offset: opts.Offset, Count: len(prods)}}, nil
}
