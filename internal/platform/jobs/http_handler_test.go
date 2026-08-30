package jobs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-fuego/fuego"
)

func TestHTTPHandler_Get(t *testing.T) {
	queue := NewQueue(1)
	pool, err := NewWorkerPool(queue, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHTTPHandler(queue, pool)

	job := newTestJob(t, "job-1")
	if err := queue.Enqueue(t.Context(), job); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.MarkRunning(t.Context(), "job-1", time.Now()); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.ReportProgress(t.Context(), "job-1", 25); err != nil {
		t.Fatal(err)
	}

	s := fuego.NewServer()
	fuego.Get(s, "/jobs/{id}", handler.Get)
	req := httptest.NewRequest(http.MethodGet, "/jobs/job-1", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var response JobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.ID != "job-1" || response.Status != StatusRunning.String() || response.Progress != 25 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestHTTPHandler_GetNotFound(t *testing.T) {
	queue := NewQueue(1)
	pool, err := NewWorkerPool(queue, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHTTPHandler(queue, pool)

	s := fuego.NewServer()
	fuego.Get(s, "/jobs/{id}", handler.Get)
	req := httptest.NewRequest(http.MethodGet, "/jobs/missing", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHTTPHandler_Cancel(t *testing.T) {
	queue := NewQueue(1)
	pool, err := NewWorkerPool(queue, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHTTPHandler(queue, pool)

	if err := queue.Enqueue(t.Context(), newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	s := fuego.NewServer()
	fuego.Put(s, "/jobs/{id}/cancel", handler.Cancel)
	req := httptest.NewRequest(http.MethodPut, "/jobs/job-1/cancel", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var response JobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Status != StatusCancelled.String() || !response.CancelRequested {
		t.Fatalf("expected cancelled response, got %+v", response)
	}
}
