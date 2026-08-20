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

// ProductStatisticsResponse is the JSON response for the product statistics report.
type ProductStatisticsResponse struct {
	Rows []ProductStatisticsRowResponse `json:"rows"`
}

// ProductStatisticsRowResponse is one row in the product statistics report response.
type ProductStatisticsRowResponse struct {
	ProductSKU    string `json:"productSku"`
	ProductName   string `json:"productName"`
	TotalQuantity int    `json:"totalQuantity"`
	EntryCount    int    `json:"entryCount"`
	EmployeeCount int    `json:"employeeCount"`
}

// DailyEmployeeProductionResponse is the JSON response for the daily employee production report.
type DailyEmployeeProductionResponse struct {
	Rows []DailyEmployeeProductionRowResponse `json:"rows"`
}

// DailyEmployeeProductionRowResponse is one row in the daily employee production report response.
type DailyEmployeeProductionRowResponse struct {
	Day           time.Time `json:"day"`
	ProductSKU    string    `json:"productSku"`
	ProductName   string    `json:"productName"`
	EmployeeID    string    `json:"employeeId"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	TotalQuantity int       `json:"totalQuantity"`
	EntryCount    int       `json:"entryCount"`
}

// EmployeeProductivityProductsResponse is the JSON response for employee productivity by product.
type EmployeeProductivityProductsResponse struct {
	Rows []EmployeeProductivityProductRowResponse `json:"rows"`
}

// EmployeeProductivityProductRowResponse is one row in the employee productivity by product response.
type EmployeeProductivityProductRowResponse struct {
	EmployeeID    string `json:"employeeId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	ProductSKU    string `json:"productSku"`
	ProductName   string `json:"productName"`
	TotalQuantity int    `json:"totalQuantity"`
	EntryCount    int    `json:"entryCount"`
}

// DailyProduction handles GET /reports/daily-production.
// Reports daily quantities for each Product (at given period).
// Answers question: "How much of each product was made each day?".
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
// Reports employee productivity for all products (at given period).
// Answers question: "How much products did each employee made overall?".
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

// ProductStatistics handles GET /reports/product-statistics.
// Reports total quantities for each Product (at given period).
// Answers question: "How much of each product was made?".
func (h *Handler) ProductStatistics(c fuego.ContextNoBody) (ProductStatisticsResponse, error) {
	from, to, err := reportRange(c)
	if err != nil {
		return ProductStatisticsResponse{}, invalidRangeError(err)
	}

	rows, err := h.store.ProductStatistics(c.Context(), from, to)
	if err != nil {
		if errors.Is(err, ErrInvalidRange) {
			return ProductStatisticsResponse{}, invalidRangeError(err)
		}
		return ProductStatisticsResponse{}, err
	}

	return ProductStatisticsResponse{Rows: productStatisticsRowsResponse(rows)}, nil
}

// DailyEmployeeProduction handles GET /reports/daily-employee-production.
// Reports daily quantities for each product and employee (at given period).
// Answers question: "Who produced how much of each product each day?".
func (h *Handler) DailyEmployeeProduction(c fuego.ContextNoBody) (DailyEmployeeProductionResponse, error) {
	from, to, err := reportRange(c)
	if err != nil {
		return DailyEmployeeProductionResponse{}, invalidRangeError(err)
	}

	rows, err := h.store.DailyEmployeeProduction(c.Context(), from, to)
	if err != nil {
		if errors.Is(err, ErrInvalidRange) {
			return DailyEmployeeProductionResponse{}, invalidRangeError(err)
		}
		return DailyEmployeeProductionResponse{}, err
	}

	return DailyEmployeeProductionResponse{Rows: dailyEmployeeProductionRowsResponse(rows)}, nil
}

// EmployeeProductivityProducts handles GET /reports/employee-productivity/products.
// Reports employee productivity broken down by product (at given period).
// Answers question: "Which products did each employee produce?".
func (h *Handler) EmployeeProductivityProducts(c fuego.ContextNoBody) (EmployeeProductivityProductsResponse, error) {
	from, to, err := reportRange(c)
	if err != nil {
		return EmployeeProductivityProductsResponse{}, invalidRangeError(err)
	}

	rows, err := h.store.EmployeeProductivityProducts(c.Context(), from, to)
	if err != nil {
		if errors.Is(err, ErrInvalidRange) {
			return EmployeeProductivityProductsResponse{}, invalidRangeError(err)
		}
		return EmployeeProductivityProductsResponse{}, err
	}

	return EmployeeProductivityProductsResponse{Rows: employeeProductivityProductRowsResponse(rows)}, nil
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

func productStatisticsRowsResponse(rows []ProductStatisticsRow) []ProductStatisticsRowResponse {
	response := make([]ProductStatisticsRowResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, ProductStatisticsRowResponse(row))
	}
	return response
}

func dailyEmployeeProductionRowsResponse(rows []DailyEmployeeProductionRow) []DailyEmployeeProductionRowResponse {
	response := make([]DailyEmployeeProductionRowResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, DailyEmployeeProductionRowResponse{
			Day:           row.Day.UTC(),
			ProductSKU:    row.ProductSKU,
			ProductName:   row.ProductName,
			EmployeeID:    row.EmployeeID,
			FirstName:     row.FirstName,
			LastName:      row.LastName,
			TotalQuantity: row.TotalQuantity,
			EntryCount:    row.EntryCount,
		})
	}
	return response
}

func employeeProductivityProductRowsResponse(rows []EmployeeProductivityProductRow) []EmployeeProductivityProductRowResponse {
	response := make([]EmployeeProductivityProductRowResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, EmployeeProductivityProductRowResponse(row))
	}
	return response
}

func invalidRangeError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{Err: err, Detail: err.Error()}
}
