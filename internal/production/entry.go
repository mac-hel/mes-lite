package production

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// ErrNotFound is returned when a production entry cannot be found by ID.
var ErrNotFound = errors.New("production entry not found")

// ErrAlreadyExists is returned when trying to create a production entry with a duplicate ID.
var ErrAlreadyExists = errors.New("production entry already exists")

// ErrInvalidEntry is returned when production entry data breaks domain rules.
var ErrInvalidEntry = errors.New("invalid production entry")

// ErrRequestConflict is returned when a request ID is reused for different production data.
var ErrRequestConflict = errors.New("production request id conflict")

// Entry records completed work for one employee and product.
type Entry struct {
	ID          string
	RequestID   string
	EmployeeID  string
	ProductSKU  string
	Quantity    int
	Workstation string
	Timestamp   time.Time
	Comment     string
}

// NewEntry creates a valid production entry and normalizes text fields.
func NewEntry(id, employeeID, productSKU string, quantity int, workstation string, timestamp time.Time, comment string) (Entry, error) {
	return NewEntryWithRequestID(id, "", employeeID, productSKU, quantity, workstation, timestamp, comment)
}

// NewEntryWithRequestID creates a valid production entry with an optional idempotency request ID.
func NewEntryWithRequestID(id, requestID, employeeID, productSKU string, quantity int, workstation string, timestamp time.Time, comment string) (Entry, error) {
	entry := Entry{
		ID:          strings.TrimSpace(id),
		RequestID:   strings.TrimSpace(requestID),
		EmployeeID:  strings.TrimSpace(employeeID),
		ProductSKU:  strings.TrimSpace(productSKU),
		Quantity:    quantity,
		Workstation: strings.TrimSpace(workstation),
		Timestamp:   timestamp.UTC(),
		Comment:     strings.TrimSpace(comment),
	}
	if err := entry.Validate(); err != nil {
		return Entry{}, err
	}

	return entry, nil
}

// NewEntryID creates a UUID-shaped identifier using only the standard library.
func NewEntryID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate production entry id: %w", err)
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	var dst [36]byte
	hex.Encode(dst[0:8], b[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], b[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], b[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], b[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], b[10:16])

	return string(dst[:]), nil
}

// Validate checks the production entry invariants that must hold in every entry point.
func (e Entry) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidEntry)
	}
	if len(strings.TrimSpace(e.RequestID)) > 128 {
		return fmt.Errorf("request id must be at most 128 characters: %w", ErrInvalidEntry)
	}
	if strings.TrimSpace(e.EmployeeID) == "" {
		return fmt.Errorf("employee id is required: %w", ErrInvalidEntry)
	}
	if strings.TrimSpace(e.ProductSKU) == "" {
		return fmt.Errorf("product sku is required: %w", ErrInvalidEntry)
	}
	if e.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero: %w", ErrInvalidEntry)
	}
	if e.Quantity > math.MaxInt32 {
		return fmt.Errorf("quantity must fit PostgreSQL integer: %w", ErrInvalidEntry)
	}
	if strings.TrimSpace(e.Workstation) == "" {
		return fmt.Errorf("workstation is required: %w", ErrInvalidEntry)
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required: %w", ErrInvalidEntry)
	}

	return nil
}
