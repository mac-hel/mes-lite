package production

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	defaultListLimit = 50
	maxListLimit     = 100
)

// ErrInvalidListOptions is returned when production-entry review filters are invalid.
var ErrInvalidListOptions = errors.New("invalid production entry list options")

// ListOptions defines filtering and pagination for production-entry review.
type ListOptions struct {
	Limit       int
	Offset      int
	EmployeeID  string
	ProductSKU  string
	Workstation string
	From        time.Time
	To          time.Time
}

// Page describes the returned production-entry page.
type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

func (o ListOptions) normalize() (ListOptions, error) {
	o.EmployeeID = strings.TrimSpace(o.EmployeeID)
	o.ProductSKU = strings.TrimSpace(o.ProductSKU)
	o.Workstation = strings.TrimSpace(o.Workstation)
	if !o.From.IsZero() {
		o.From = o.From.UTC()
	}
	if !o.To.IsZero() {
		o.To = o.To.UTC()
	}

	if o.Limit == 0 {
		o.Limit = defaultListLimit
	}
	if o.Limit < 0 || o.Limit > maxListLimit {
		return ListOptions{}, fmt.Errorf("limit must be between 1 and %d: %w", maxListLimit, ErrInvalidListOptions)
	}
	if o.Offset < 0 {
		return ListOptions{}, fmt.Errorf("offset must be greater than or equal to zero: %w", ErrInvalidListOptions)
	}
	if !o.From.IsZero() && !o.To.IsZero() && !o.From.Before(o.To) {
		return ListOptions{}, fmt.Errorf("from must be before to: %w", ErrInvalidListOptions)
	}

	return o, nil
}
