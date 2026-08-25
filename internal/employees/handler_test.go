package employees

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-fuego/fuego"
)

func TestHandler_Create(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)

	body := `{"id":"001","firstName":"John","lastName":"Doe","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	responseBody := w.Body.Bytes()
	var emp EmployeeResponse
	if err := json.Unmarshal(responseBody, &emp); err != nil {
		t.Fatal(err)
	}

	if emp.ID != "001" {
		t.Errorf("expected ID 001, got %q", emp.ID)
	}
	if emp.FirstName != "John" {
		t.Errorf("expected FirstName John, got %q", emp.FirstName)
	}
	if emp.IsActive != true {
		t.Errorf("expected IsActive true, got %v", emp.IsActive)
	}

	var fields map[string]any
	if err := json.Unmarshal(responseBody, &fields); err != nil {
		t.Fatal(err)
	}
	if _, ok := fields["firstName"]; !ok {
		t.Fatal("expected lower-camel firstName field")
	}
	if _, ok := fields["FirstName"]; ok {
		t.Fatal("did not expect capitalized FirstName field")
	}
}

func TestHandler_DuplicateCreate(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)

	body := `{"id":"001","firstName":"John","lastName":"Doe","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Create_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"missing firstName", `{"id":"001","lastName":"Doe","email":"a@b.com"}`, http.StatusBadRequest},
		{"missing lastName", `{"id":"001","firstName":"John","email":"a@b.com"}`, http.StatusBadRequest},
		{"missing email", `{"id":"001","firstName":"John","lastName":"Doe"}`, http.StatusBadRequest},
		{"missing id", `{"firstName":"John","lastName":"Doe","email":"a@b.com"}`, http.StatusBadRequest},
		{"invalid email", `{"id":"001","firstName":"John","lastName":"Doe","email":"notanemail"}`, http.StatusBadRequest},
	}

	store := NewInMemoryStore()
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_List_Empty(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Get(s, "/employees", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/employees", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListEmployeesResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Employees == nil {
		t.Error("expected non-nil slice for empty list")
	}
	if len(resp.Employees) != 0 {
		t.Errorf("expected 0 employees, got %d", len(resp.Employees))
	}
}

func TestHandler_List_AfterCreate(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)
	fuego.Get(s, "/employees", handler.List)

	createEmployee := func(id string) {
		body := `{"id":"` + id + `","firstName":"Test","lastName":"User","email":"` + id + `@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("create %s: expected 200, got %d", id, w.Code)
		}
	}

	createEmployee("001")
	createEmployee("002")
	createEmployee("003")

	req := httptest.NewRequest(http.MethodGet, "/employees", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListEmployeesResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Employees) != 3 {
		t.Fatalf("expected 3 employees, got %d", len(resp.Employees))
	}

	ids := make(map[string]bool, 3)
	for _, emp := range resp.Employees {
		ids[emp.ID] = true
	}
	for _, expected := range []string{"001", "002", "003"} {
		if !ids[expected] {
			t.Errorf("expected employee %s in list", expected)
		}
	}
}

func TestHandler_List_QueryOptions(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)
	fuego.Get(s, "/employees", handler.List)
	fuego.Put(s, "/employees/{id}/deactivate", handler.Deactivate)

	createEmployee := func(id, firstName, lastName string) {
		body := `{"id":"` + id + `","firstName":"` + firstName + `","lastName":"` + lastName + `","email":"` + id + `@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("create %s: expected 200, got %d", id, w.Code)
		}
	}

	createEmployee("001", "Ana", "Alpha")
	createEmployee("002", "Bob", "Beta")
	createEmployee("003", "Carla", "Zulu")

	req := httptest.NewRequest(http.MethodPut, "/employees/002/deactivate", nil)
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodGet, "/employees?active=true&sort=-name&limit=1&offset=1&q=a", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp ListEmployeesResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(resp.Employees))
	}
	if resp.Employees[0].ID != "001" {
		t.Fatalf("expected paginated employee 001, got %q", resp.Employees[0].ID)
	}
	if resp.Pagination.Limit != 1 || resp.Pagination.Offset != 1 || resp.Pagination.Count != 1 {
		t.Fatalf("unexpected pagination: %#v", resp.Pagination)
	}
}

func TestHandler_List_InvalidQueryOptions(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/employees", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/employees?sort=unknown", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Update(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)
	fuego.Put(s, "/employees/{id}", handler.Update)

	createBody := `{"id":"001","firstName":"John","lastName":"Doe","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	updateBody := `{"firstName":"Jane","lastName":"Smith","email":"jane@example.com","version":1}`
	req = httptest.NewRequest(http.MethodPut, "/employees/001", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var emp EmployeeResponse
	if err := json.NewDecoder(w.Body).Decode(&emp); err != nil {
		t.Fatal(err)
	}

	if emp.ID != "001" {
		t.Errorf("expected ID 001, got %q", emp.ID)
	}
	if emp.FirstName != "Jane" {
		t.Errorf("expected FirstName Jane, got %q", emp.FirstName)
	}
	if emp.LastName != "Smith" {
		t.Errorf("expected LastName Smith, got %q", emp.LastName)
	}
	if emp.Email != "jane@example.com" {
		t.Errorf("expected Email jane@example.com, got %q", emp.Email)
	}
	if emp.Version != 2 {
		t.Errorf("expected Version 2, got %d", emp.Version)
	}
}

func TestHandler_Update_VersionConflict(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)
	fuego.Put(s, "/employees/{id}", handler.Update)

	createBody := `{"id":"001","firstName":"John","lastName":"Doe","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	updateBody := `{"firstName":"Jane","lastName":"Smith","email":"jane@example.com","version":1}`
	req = httptest.NewRequest(http.MethodPut, "/employees/001", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodPut, "/employees/001", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Update_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Put(s, "/employees/{id}", handler.Update)

	body := `{"firstName":"Jane","lastName":"Smith","email":"jane@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/employees/nonexistent", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Deactivate(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/employees", handler.Create)
	fuego.Put(s, "/employees/{id}/deactivate", handler.Deactivate)

	createBody := `{"id":"001","firstName":"John","lastName":"Doe","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodPut, "/employees/001/deactivate", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var emp EmployeeResponse
	if err := json.NewDecoder(w.Body).Decode(&emp); err != nil {
		t.Fatal(err)
	}

	if emp.ID != "001" {
		t.Errorf("expected ID 001, got %q", emp.ID)
	}
	if emp.IsActive != false {
		t.Errorf("expected IsActive false, got %v", emp.IsActive)
	}
}

func TestHandler_Deactivate_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Put(s, "/employees/{id}/deactivate", handler.Deactivate)

	req := httptest.NewRequest(http.MethodPut, "/employees/nonexistent/deactivate", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}
