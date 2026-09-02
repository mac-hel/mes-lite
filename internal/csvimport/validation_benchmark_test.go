package csvimport

import (
	"fmt"
	"strings"
	"testing"
)

var benchmarkValidationResult ValidationResult

func BenchmarkValidateProductionEntries(b *testing.B) {
	benchmarks := []struct {
		name string
		rows int
	}{
		{name: "100_rows", rows: 100},
		{name: "1000_rows", rows: 1000},
		{name: "10000_rows", rows: 10000},
	}

	for _, bm := range benchmarks {
		input := productionEntryBenchmarkCSV(bm.rows)

		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()

			for range b.N {
				reader, err := NewProductionEntryReader(strings.NewReader(input))
				if err != nil {
					b.Fatalf("NewProductionEntryReader() error = %v", err)
				}

				result, err := ValidateProductionEntries(reader)
				if err != nil {
					b.Fatalf("ValidateProductionEntries() error = %v", err)
				}
				benchmarkValidationResult = result
			}
		})
	}
}

func productionEntryBenchmarkCSV(rows int) string {
	var builder strings.Builder
	builder.Grow(80 * (rows + 1))
	builder.WriteString("employee_id,product_sku,quantity,workstation,timestamp,comment\n")

	for i := range rows {
		_, _ = fmt.Fprintf(
			&builder,
			"emp-%d,sku-%d,%d,ws-%d,2026-08-20T10:00:00Z,benchmark row %d\n",
			i%50,
			i%25,
			(i%100)+1,
			i%10,
			i,
		)
	}

	return builder.String()
}
