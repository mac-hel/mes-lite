package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-fuego/fuego"
)

// Canceller requests cancellation for a job.
type Canceller interface {
	Cancel(ctx context.Context, id string) (Job, error)
}

// HTTPHandler exposes job status operations over HTTP.
type HTTPHandler struct {
	queue     *Queue
	canceller Canceller
}

// NewHTTPHandler creates a background-job HTTP handler.
func NewHTTPHandler(queue *Queue, canceller Canceller) *HTTPHandler {
	return &HTTPHandler{queue: queue, canceller: canceller}
}

// JobResponse is the HTTP representation of a background job.
type JobResponse struct {
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Progress        int       `json:"progress"`
	CancelRequested bool      `json:"cancelRequested"`
	EnqueuedAt      time.Time `json:"enqueuedAt"`
	StartedAt       time.Time `json:"startedAt,omitempty"`
	FinishedAt      time.Time `json:"finishedAt,omitempty"`
	Error           string    `json:"error,omitempty"`
	Result          any       `json:"result,omitempty"`
}

// Get returns one background job status snapshot.
func (h *HTTPHandler) Get(c fuego.ContextNoBody) (JobResponse, error) {
	job, err := h.queue.Find(c.Context(), c.PathParam("id"))
	if err != nil {
		return JobResponse{}, jobHTTPError(c.PathParam("id"), err)
	}
	return NewJobResponse(job), nil
}

// Cancel requests cancellation for one background job.
func (h *HTTPHandler) Cancel(c fuego.ContextNoBody) (JobResponse, error) {
	job, err := h.canceller.Cancel(c.Context(), c.PathParam("id"))
	if err != nil {
		return JobResponse{}, jobHTTPError(c.PathParam("id"), err)
	}
	return NewJobResponse(job), nil
}

// NewJobResponse maps a tracked job to its HTTP response shape.
func NewJobResponse(job Job) JobResponse {
	response := JobResponse{
		ID:              job.ID,
		Type:            job.Type.String(),
		Status:          job.Status.String(),
		Progress:        job.Progress,
		CancelRequested: job.CancelRequested,
		EnqueuedAt:      job.EnqueuedAt,
		StartedAt:       job.StartedAt,
		FinishedAt:      job.FinishedAt,
		Error:           job.Error,
	}
	if len(job.Result) > 0 {
		response.Result = json.RawMessage(job.Result)
	}
	return response
}

func jobHTTPError(id string, err error) error {
	if errors.Is(err, ErrNotFound) {
		return fuego.NotFoundError{Err: err, Detail: fmt.Sprintf("background job %q not found", id)}
	}
	return err
}
