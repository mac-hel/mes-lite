// Package jobs models background work that should not block an HTTP request.
package jobs

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrInvalidJob is returned when background job data breaks domain rules.
var ErrInvalidJob = errors.New("invalid background job")

// maxTypeLength bounds the job type so a malformed value cannot become an unbounded map key.
const maxTypeLength = 64

// Status describes where a background job is in its lifecycle.
type Status string

// Background job lifecycle states.
const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// String returns the wire representation of the status.
func (s Status) String() string {
	return string(s)
}

// Valid reports whether the status is one of the known lifecycle states.
func (s Status) Valid() bool {
	switch s {
	case StatusQueued, StatusRunning, StatusSucceeded, StatusFailed, StatusCancelled:
		return true
	default:
		return false
	}
}

// Terminal reports whether a job in this status will never run again.
func (s Status) Terminal() bool {
	switch s {
	case StatusSucceeded, StatusFailed, StatusCancelled:
		return true
	default:
		return false
	}
}

// Type identifies which background workload a job carries.
type Type string

// TypeProductionEntryImport imports historical production entries from CSV data.
const TypeProductionEntryImport Type = "production-entry-import"

// String returns the wire representation of the job type.
func (t Type) String() string {
	return string(t)
}

// Job is one unit of background work tracked by the application.
//
// Job is a value type. The queue owns the canonical copy and hands out copies,
// so callers can never mutate tracked job state by accident.
type Job struct {
	ID         string
	Type       Type
	Status     Status
	Payload    []byte
	EnqueuedAt time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	Error      string
}

// NewJob creates a queued background job and copies the payload so the caller
// cannot mutate job state through the slice it passed in.
func NewJob(id string, jobType Type, payload []byte, enqueuedAt time.Time) (Job, error) {
	job := Job{
		ID:         strings.TrimSpace(id),
		Type:       Type(strings.TrimSpace(string(jobType))),
		Status:     StatusQueued,
		Payload:    copyPayload(payload),
		EnqueuedAt: enqueuedAt.UTC(),
	}
	if err := job.Validate(); err != nil {
		return Job{}, err
	}

	return job, nil
}

// NewJobID creates a UUID-shaped identifier using only the standard library.
func NewJobID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate background job id: %w", err)
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

// Validate checks the job invariants that must hold wherever a job is created.
func (j Job) Validate() error {
	if strings.TrimSpace(j.ID) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidJob)
	}
	if strings.TrimSpace(string(j.Type)) == "" {
		return fmt.Errorf("type is required: %w", ErrInvalidJob)
	}
	if len(j.Type) > maxTypeLength {
		return fmt.Errorf("type must be at most %d characters: %w", maxTypeLength, ErrInvalidJob)
	}
	if !j.Status.Valid() {
		return fmt.Errorf("status %q is not supported: %w", j.Status, ErrInvalidJob)
	}
	if j.EnqueuedAt.IsZero() {
		return fmt.Errorf("enqueued at is required: %w", ErrInvalidJob)
	}

	return nil
}

// clone returns a copy that shares no memory with the receiver.
func (j Job) clone() Job {
	j.Payload = copyPayload(j.Payload)
	return j
}

// copyPayload returns a defensive copy so a payload slice is never shared
// between the caller and the queue.
func copyPayload(payload []byte) []byte {
	if payload == nil {
		return nil
	}

	copied := make([]byte, len(payload))
	copy(copied, payload)
	return copied
}
