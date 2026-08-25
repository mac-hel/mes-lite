package production

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrInvalidCorrection is returned when a correction breaks audit rules.
var ErrInvalidCorrection = errors.New("invalid production entry correction")

// Correction records an append-only correction for a production entry.
type Correction struct {
	ID          string    `json:"id"`
	EntryID     string    `json:"entryId"`
	ActorUserID string    `json:"actorUserId"`
	Reason      string    `json:"reason"`
	EmployeeID  string    `json:"employeeId"`
	ProductSKU  string    `json:"productSku"`
	Quantity    int       `json:"quantity"`
	Workstation string    `json:"workstation"`
	Timestamp   time.Time `json:"timestamp"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"createdAt"`
}

// NewCorrection creates a valid append-only correction and normalizes fields.
func NewCorrection(id, entryID, actorUserID, reason, employeeID, productSKU string, quantity int, workstation string, timestamp time.Time, comment string) (Correction, error) {
	correction := Correction{
		ID:          strings.TrimSpace(id),
		EntryID:     strings.TrimSpace(entryID),
		ActorUserID: strings.TrimSpace(actorUserID),
		Reason:      strings.TrimSpace(reason),
		EmployeeID:  strings.TrimSpace(employeeID),
		ProductSKU:  strings.TrimSpace(productSKU),
		Quantity:    quantity,
		Workstation: strings.TrimSpace(workstation),
		Timestamp:   timestamp.UTC(),
		Comment:     strings.TrimSpace(comment),
	}
	if err := correction.Validate(); err != nil {
		return Correction{}, err
	}

	return correction, nil
}

// Validate checks correction invariants.
func (c Correction) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidCorrection)
	}
	if strings.TrimSpace(c.EntryID) == "" {
		return fmt.Errorf("entry id is required: %w", ErrInvalidCorrection)
	}
	if strings.TrimSpace(c.ActorUserID) == "" {
		return fmt.Errorf("actor user id is required: %w", ErrInvalidCorrection)
	}
	if strings.TrimSpace(c.Reason) == "" {
		return fmt.Errorf("reason is required: %w", ErrInvalidCorrection)
	}
	entry, err := NewEntry("00000000-0000-4000-8000-000000000000", c.EmployeeID, c.ProductSKU, c.Quantity, c.Workstation, c.Timestamp, c.Comment)
	if err != nil {
		return fmt.Errorf("corrected entry values: %w", ErrInvalidCorrection)
	}
	if entry.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required: %w", ErrInvalidCorrection)
	}

	return nil
}
