package reporting

import (
	"errors"
	"testing"
	"time"
)

func TestValidateRange(t *testing.T) {
	from := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	tests := []struct {
		name    string
		from    time.Time
		to      time.Time
		wantErr bool
	}{
		{name: "valid", from: from, to: to},
		{name: "missing from", to: to, wantErr: true},
		{name: "missing to", from: from, wantErr: true},
		{name: "same instant", from: from, to: from, wantErr: true},
		{name: "from after to", from: to, to: from, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRange(tt.from, tt.to)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidRange) {
					t.Fatalf("validateRange() error = %v, want ErrInvalidRange", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateRange() error = %v, want nil", err)
			}
		})
	}
}
