package csvimport

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestProductionEntryReader_ReadsRowsSequentially(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		" emp-1 , sku-1 , 12 , ws-1 , 2026-08-20T10:00:00Z , first batch ",
		"emp-2,sku-2,7,ws-2,2026-08-20T11:00:00Z,second batch",
	}, "\n"))

	reader, err := NewProductionEntryReader(input)
	if err != nil {
		t.Fatalf("NewProductionEntryReader() error = %v", err)
	}

	first, err := reader.Read()
	if err != nil {
		t.Fatalf("Read() first row error = %v", err)
	}
	if first.RowNumber != 2 {
		t.Fatalf("first row number = %d, want 2", first.RowNumber)
	}
	if first.EmployeeID != "emp-1" || first.ProductSKU != "sku-1" || first.Quantity != "12" || first.Workstation != "ws-1" || first.Timestamp != "2026-08-20T10:00:00Z" || first.Comment != "first batch" {
		t.Fatalf("first row = %+v", first)
	}

	second, err := reader.Read()
	if err != nil {
		t.Fatalf("Read() second row error = %v", err)
	}
	if second.RowNumber != 3 {
		t.Fatalf("second row number = %d, want 3", second.RowNumber)
	}
	if second.EmployeeID != "emp-2" || second.ProductSKU != "sku-2" || second.Quantity != "7" || second.Workstation != "ws-2" || second.Timestamp != "2026-08-20T11:00:00Z" || second.Comment != "second batch" {
		t.Fatalf("second row = %+v", second)
	}

	_, err = reader.Read()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Read() exhausted error = %v, want io.EOF", err)
	}
}

func TestProductionEntryReader_AcceptsHeaderCaseAndSpace(t *testing.T) {
	input := strings.NewReader(" Employee_ID , Product_SKU , Quantity , Workstation , Timestamp , Comment \n")

	_, err := NewProductionEntryReader(input)
	if err != nil {
		t.Fatalf("NewProductionEntryReader() error = %v", err)
	}
}

func TestProductionEntryReader_RejectsMissingHeader(t *testing.T) {
	_, err := NewProductionEntryReader(strings.NewReader(""))
	if !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("NewProductionEntryReader() error = %v, want ErrInvalidHeader", err)
	}
}

func TestProductionEntryReader_RejectsUnexpectedHeader(t *testing.T) {
	input := strings.NewReader("employee_id,product_sku,quantity\n")

	_, err := NewProductionEntryReader(input)
	if !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("NewProductionEntryReader() error = %v, want ErrInvalidHeader", err)
	}
}

func TestProductionEntryReader_ReturnsCSVParseError(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,\"unterminated",
	}, "\n"))

	reader, err := NewProductionEntryReader(input)
	if err != nil {
		t.Fatalf("NewProductionEntryReader() error = %v", err)
	}

	_, err = reader.Read()
	if err == nil {
		t.Fatal("Read() error = nil, want CSV parse error")
	}
	if errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("Read() error = %v, did not want ErrInvalidRecord", err)
	}
}

func TestProductionEntryReader_DoesNotReadPastHeaderDuringConstruction(t *testing.T) {
	input := &singleRecordReader{
		data: "employee_id,product_sku,quantity,workstation,timestamp,comment\nemp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,done\n",
	}

	reader, err := NewProductionEntryReader(input)
	if err != nil {
		t.Fatalf("NewProductionEntryReader() error = %v", err)
	}
	if input.recordsRead != 1 {
		t.Fatalf("records read after constructor = %d, want 1", input.recordsRead)
	}

	_, err = reader.Read()
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if input.recordsRead != 2 {
		t.Fatalf("records read after first row = %d, want 2", input.recordsRead)
	}
}

type singleRecordReader struct {
	data        string
	recordsRead int
}

func (r *singleRecordReader) Read(p []byte) (int, error) {
	if r.data == "" {
		return 0, io.EOF
	}

	newline := strings.IndexByte(r.data, '\n')
	if newline == -1 {
		newline = len(r.data) - 1
	}
	chunk := r.data[:newline+1]
	r.data = r.data[newline+1:]
	r.recordsRead++

	return copy(p, chunk), nil
}
