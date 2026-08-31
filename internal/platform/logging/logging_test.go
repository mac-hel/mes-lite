package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	logger, err := New(&buf, "info", "json")
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("database connected", "component", "server")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["msg"] != "database connected" {
		t.Fatalf("expected structured message, got %#v", got["msg"])
	}
	if got["component"] != "server" {
		t.Fatalf("expected structured component, got %#v", got["component"])
	}
}

func TestNewTextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger, err := New(&buf, "debug", "text")
	if err != nil {
		t.Fatal(err)
	}

	logger.Debug("debug enabled", "component", "test")

	out := buf.String()
	if !strings.Contains(out, "msg=\"debug enabled\"") {
		t.Fatalf("expected text slog output, got %q", out)
	}
	if !strings.Contains(out, "component=test") {
		t.Fatalf("expected structured component, got %q", out)
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name string
		want slog.Level
	}{
		{name: "", want: slog.LevelInfo},
		{name: "debug", want: slog.LevelDebug},
		{name: "info", want: slog.LevelInfo},
		{name: "warn", want: slog.LevelWarn},
		{name: "warning", want: slog.LevelWarn},
		{name: "error", want: slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLevel(tt.name)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestInvalidConfiguration(t *testing.T) {
	if _, err := New(&bytes.Buffer{}, "verbose", "json"); err == nil {
		t.Fatal("expected invalid level error")
	}
	if _, err := New(&bytes.Buffer{}, "info", "xml"); err == nil {
		t.Fatal("expected invalid format error")
	}
}
