package employees

import (
	"errors"
	"fmt"
	"strings"
)

const (
	defaultListLimit = 50
	maxListLimit     = 100
)

// ErrInvalidListOptions is returned when employee list query options are not supported.
var ErrInvalidListOptions = errors.New("invalid employee list options")

// ListOptions defines filtering, sorting and pagination for employee lists.
type ListOptions struct {
	Limit  int
	Offset int
	Sort   string
	Query  string
	Active *bool
}

// Page describes the returned page.
type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

func (o ListOptions) normalize() (ListOptions, error) {
	o.Query = strings.TrimSpace(o.Query)
	o.Sort = strings.TrimSpace(o.Sort)

	if o.Limit == 0 {
		o.Limit = defaultListLimit
	}
	if o.Limit < 0 || o.Limit > maxListLimit {
		return ListOptions{}, fmt.Errorf("limit must be between 1 and %d: %w", maxListLimit, ErrInvalidListOptions)
	}
	if o.Offset < 0 {
		return ListOptions{}, fmt.Errorf("offset must be greater than or equal to zero: %w", ErrInvalidListOptions)
	}
	if o.Sort == "" {
		o.Sort = "id"
	}
	if !validSort(o.Sort) {
		return ListOptions{}, fmt.Errorf("sort %q is not supported: %w", o.Sort, ErrInvalidListOptions)
	}

	return o, nil
}

func validSort(sort string) bool {
	switch sort {
	case "id", "-id", "name", "-name", "email", "-email":
		return true
	default:
		return false
	}
}

func activeFilter(active *bool) string {
	if active == nil {
		return ""
	}
	if *active {
		return "true"
	}
	return "false"
}
