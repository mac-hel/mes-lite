package jobs

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewJob_NormalizesAndQueues(t *testing.T) {
	enqueuedAt := time.Date(2026, 8, 27, 10, 30, 0, 0, time.FixedZone("CEST", 2*60*60))

	job, err := NewJob(" job-1 ", Type(" "+string(TypeProductionEntryImport)+" "), []byte("sku,qty"), enqueuedAt)
	if err != nil {
		t.Fatal(err)
	}

	if job.ID != "job-1" {
		t.Errorf("expected trimmed ID, got %q", job.ID)
	}
	if job.Type != TypeProductionEntryImport {
		t.Errorf("expected trimmed type, got %q", job.Type)
	}
	if job.Status != StatusQueued {
		t.Errorf("expected new job to be queued, got %q", job.Status)
	}
	if job.EnqueuedAt.Location() != time.UTC {
		t.Errorf("expected UTC enqueued at, got %s", job.EnqueuedAt.Location())
	}
	if !job.StartedAt.IsZero() || !job.FinishedAt.IsZero() {
		t.Error("expected a new job to have no start or finish timestamp")
	}
}

func TestNewJob_CopiesPayload(t *testing.T) {
	payload := []byte("original")

	job, err := NewJob("job-1", TypeProductionEntryImport, payload, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	payload[0] = 'X'

	if string(job.Payload) != "original" {
		t.Errorf("expected the job to own its payload copy, got %q", job.Payload)
	}
}

func TestJob_Validate(t *testing.T) {
	valid, err := NewJob("job-1", TypeProductionEntryImport, nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		job     Job
		wantErr bool
	}{
		{"valid", valid, false},
		{"missing id", Job{Type: TypeProductionEntryImport, Status: StatusQueued, EnqueuedAt: time.Now()}, true},
		{"missing type", Job{ID: "job-1", Status: StatusQueued, EnqueuedAt: time.Now()}, true},
		{"type too long", Job{ID: "job-1", Type: Type(strings.Repeat("t", maxTypeLength+1)), Status: StatusQueued, EnqueuedAt: time.Now()}, true},
		{"unknown status", Job{ID: "job-1", Type: TypeProductionEntryImport, Status: Status("archived"), EnqueuedAt: time.Now()}, true},
		{"missing enqueued at", Job{ID: "job-1", Type: TypeProductionEntryImport, Status: StatusQueued}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.job.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidJob) {
					t.Fatalf("expected ErrInvalidJob, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected valid job, got %v", err)
			}
		})
	}
}

func TestStatus_ValidAndTerminal(t *testing.T) {
	tests := []struct {
		status       Status
		wantValid    bool
		wantTerminal bool
	}{
		{StatusQueued, true, false},
		{StatusRunning, true, false},
		{StatusSucceeded, true, true},
		{StatusFailed, true, true},
		{StatusCancelled, true, true},
		{Status("archived"), false, false},
		{Status(""), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.Valid(); got != tt.wantValid {
				t.Errorf("Valid() = %v, want %v", got, tt.wantValid)
			}
			if got := tt.status.Terminal(); got != tt.wantTerminal {
				t.Errorf("Terminal() = %v, want %v", got, tt.wantTerminal)
			}
		})
	}
}

func TestNewJobID_IsUnpredictableAndUUIDShaped(t *testing.T) {
	first, err := NewJobID()
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewJobID()
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Error("expected unique job IDs")
	}
	if len(first) != 36 {
		t.Errorf("expected a UUID-shaped ID, got %q", first)
	}
}
