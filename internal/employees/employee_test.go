package employees_test

import (
	"errors"
	"testing"

	"github.com/mac-hel/mes-lite/internal/employees"
)

func TestEmployeeZeroValue(t *testing.T) {
	var emp employees.Employee

	if emp.ID != "" {
		t.Errorf("expected empty ID, got %q", emp.ID)
	}
	if emp.FirstName != "" {
		t.Errorf("expected empty FirstName, got %q", emp.FirstName)
	}
	if emp.LastName != "" {
		t.Errorf("expected empty LastName, got %q", emp.LastName)
	}
	if emp.Email != "" {
		t.Errorf("expected empty Email, got %q", emp.Email)
	}
	if emp.IsActive != false {
		t.Errorf("expected IsActive to be false, got %v", emp.IsActive)
	}
}

func TestNewEmployee(t *testing.T) {
	emp, err := employees.NewEmployee("001", "Jane", "Smith", "jane@example.com")
	if err != nil {
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
	if emp.IsActive != true {
		t.Errorf("expected IsActive to be true, got %v", emp.IsActive)
	}
}

func TestNewEmployee_RejectsInvalidState(t *testing.T) {
	_, err := employees.NewEmployee("", "Jane", "Smith", "jane@example.com")
	if !errors.Is(err, employees.ErrInvalidEmployee) {
		t.Fatalf("expected ErrInvalidEmployee, got %v", err)
	}
}

func TestEmployee_UpdateDetailsPreservesValidState(t *testing.T) {
	emp, err := employees.NewEmployee("001", "Jane", "Smith", "jane@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if err := emp.UpdateDetails(" ", "Worker", "worker@example.com"); !errors.Is(err, employees.ErrInvalidEmployee) {
		t.Fatalf("expected ErrInvalidEmployee, got %v", err)
	}
	if emp.FirstName != "Jane" {
		t.Errorf("expected failed update to preserve FirstName, got %q", emp.FirstName)
	}
}
