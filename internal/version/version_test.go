package version

import "testing"

func TestString(t *testing.T) {
	if got := String(); got != "dev" {
		t.Errorf("String() = %q, want %q", got, "dev")
	}
}

func TestFull(t *testing.T) {
	if got := Full(); got == "" {
		t.Error("Full() returned empty string")
	}
}
