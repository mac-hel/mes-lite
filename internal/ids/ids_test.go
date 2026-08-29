package ids

import (
	"strings"
	"sync"
	"testing"
)

func TestNew_CanonicalUUIDv4Format(t *testing.T) {
	id := New()

	if len(id) != 36 {
		t.Fatalf("expected 36 characters, got %d in %q", len(id), id)
	}
	for _, position := range []int{8, 13, 18, 23} {
		if id[position] != '-' {
			t.Errorf("expected a hyphen at position %d in %q", position, id)
		}
	}
	if id[14] != '4' {
		t.Errorf("expected version 4 in %q", id)
	}
	if !strings.ContainsRune("89ab", rune(id[19])) {
		t.Errorf("expected an RFC 4122 variant nibble in %q, got %q", id, id[19])
	}

	const hexDigits = "0123456789abcdef"
	for position, char := range id {
		if position == 8 || position == 13 || position == 18 || position == 23 {
			continue
		}
		if !strings.ContainsRune(hexDigits, char) {
			t.Errorf("expected lower-case hex at position %d in %q, got %q", position, id, char)
		}
	}
}

func TestNew_IsUnpredictable(t *testing.T) {
	const total = 1000

	seen := make(map[string]struct{}, total)
	for range total {
		id := New()
		if _, ok := seen[id]; ok {
			t.Fatalf("generated a duplicate id %q", id)
		}
		seen[id] = struct{}{}
	}
}

// TestNew_ConcurrentCallers is written for `go test -race`: the generator holds
// no shared state, so concurrent callers must not need external synchronization.
func TestNew_ConcurrentCallers(t *testing.T) {
	const (
		goroutines   = 8
		perGoroutine = 200
	)

	results := make(chan string, goroutines*perGoroutine)

	var group sync.WaitGroup
	for range goroutines {
		group.Add(1)
		go func() {
			defer group.Done()
			for range perGoroutine {
				results <- New()
			}
		}()
	}
	group.Wait()
	close(results)

	seen := make(map[string]struct{}, goroutines*perGoroutine)
	for id := range results {
		if _, ok := seen[id]; ok {
			t.Fatalf("generated a duplicate id %q under concurrency", id)
		}
		seen[id] = struct{}{}
	}

	if len(seen) != goroutines*perGoroutine {
		t.Fatalf("expected %d unique ids, got %d", goroutines*perGoroutine, len(seen))
	}
}
