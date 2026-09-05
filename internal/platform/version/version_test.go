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

func TestCurrent(t *testing.T) {
	info := Current()

	if info.Version != Version {
		t.Fatalf("Version = %q, want %q", info.Version, Version)
	}
	if info.Commit != Commit {
		t.Fatalf("Commit = %q, want %q", info.Commit, Commit)
	}
	if info.BuildTime != BuildTime {
		t.Fatalf("BuildTime = %q, want %q", info.BuildTime, BuildTime)
	}
}
