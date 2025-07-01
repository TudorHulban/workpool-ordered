package workpoolordered

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Simple string processor for benchmarking.
func stringProcessor(s string) (string, bool, error) {
	// Simulate some work
	return fmt.Sprintf("processed-%s", s), false, nil
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkInsert-12    	 5439858	       200.1 ns/op	      72 B/op	       2 allocs/op
func BenchmarkInsert(b *testing.B) {
	list, errCr := NewCList(
		&ParamsCList[string]{
			Processor:       stringProcessor,
			Workers:         1,
			WaitToStartWork: true,
		},
	)
	require.NoError(b, errCr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		list.Insert(fmt.Sprintf("item-%d", i))
	}
}
