package reporting

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-fuego/fuego"
)

// Handler holds HTTP handlers for reporting endpoints.
type Handler struct {
	store Store
}

// NewHandler creates a reporting handler with its read store dependency.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// DailyProductionResponse is the JSON response for the daily production report.
type DailyProductionResponse struct {
	Rows []DailyProductionRowResponse `json:"rows"`
}

// DailyProductionRowResponse is one row in the daily production report response.
type DailyProductionRowResponse struct {
	Day           time.Time `json:"day"`
	ProductSKU    string    `json:"productSku"`
	TotalQuantity int       `json:"totalQuantity"`
	EntryCount    int       `json:"entryCount"`
}

// EmployeeProductivityResponse is the JSON response for the employee productivity report.
type EmployeeProductivityResponse struct {
	Rows []EmployeeProductivityRowResponse `json:"rows"`
}

// EmployeeProductivityRowResponse is one row in the employee productivity report response.
type EmployeeProductivityRowResponse struct {
	EmployeeID    string `json:"employeeId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	TotalQuantity int    `json:"totalQuantity"`
	EntryCount    int    `json:"entryCount"`
}

// DailyProduction handles GET /reports/daily-production.
func (h *Handler) DailyProduction(c fuego.ContextNoBody) (DailyProductionResponse, error) {
	from, to, err := reportRange(c)
	if err != nil {
		return DailyProductionResponse{}, invalidRangeError(err)
	}

	rows, err := h.store.DailyProduction(c.Context(), from, to)
	if err != nil {
		if errors.Is(err, ErrInvalidRange) {
			return DailyProductionResponse{}, invalidRangeError(err)
		}
		return DailyProductionResponse{}, err
	}

	return DailyProductionResponse{Rows: dailyProductionRowsResponse(rows)}, nil
}

// EmployeeProductivity handles GET /reports/employee-productivity.
func (h *Handler) EmployeeProductivity(c fuego.ContextNoBody) (EmployeeProductivityResponse, error) {
	from, to, err := reportRange(c)
	if err != nil {
		return EmployeeProductivityResponse{}, invalidRangeError(err)
	}

	rows, err := h.store.EmployeeProductivity(c.Context(), from, to)
	if err != nil {
		if errors.Is(err, ErrInvalidRange) {
			return EmployeeProductivityResponse{}, invalidRangeError(err)
		}
		return EmployeeProductivityResponse{}, err
	}

	return EmployeeProductivityResponse{Rows: employeeProductivityRowsResponse(rows)}, nil
}

func reportRange(c fuego.ContextNoBody) (time.Time, time.Time, error) {
	from, err := parseReportTime(c.QueryParam("from"), "from")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err := parseReportTime(c.QueryParam("to"), "to")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return from, to, nil
}

func parseReportTime(raw, name string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, fmt.Errorf("%s is required: %w", name, ErrInvalidRange)
	}
	value, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be RFC3339: %w", name, ErrInvalidRange)
	}
	return value.UTC(), nil
}

func dailyProductionRowsResponse(rows []DailyProductionRow) []DailyProductionRowResponse {
	response := make([]DailyProductionRowResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, DailyProductionRowResponse{
			Day:           row.Day.UTC(),
			ProductSKU:    row.ProductSKU,
			TotalQuantity: row.TotalQuantity,
			EntryCount:    row.EntryCount,
		})
	}
	return response
}

func employeeProductivityRowsResponse(rows []EmployeeProductivityRow) []EmployeeProductivityRowResponse {
	response := make([]EmployeeProductivityRowResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, EmployeeProductivityRowResponse(row))
	}
	return response
}

func invalidRangeError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{Err: err, Detail: err.Error()}
}
