package workpoolordered

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Simple string processor for benchmarking.
func stringProcessor(s string) (string, bool, error) {
	// Simulate some work
	return fmt.Sprintf("processed-%s", s), false, nil
}

// Heavy computation processor.
func heavyProcessor(s string) (string, error) {
	// Simulate CPU-intensive work
	sum := 0

	for i := 0; i < 1000; i++ {
		sum += i
	}

	return fmt.Sprintf("heavy-%s-%d", s, sum), nil
}

// I/O simulation processor
func ioProcessor(s string) (string, error) {
	// Simulate I/O delay
	time.Sleep(time.Microsecond * 100)

	return fmt.Sprintf("io-%s", s), nil
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkInsert-16    	 9932575	       112.3 ns/op	      72 B/op	       2 allocs/op
func BenchmarkInsert(b *testing.B) {
	list, errCr := NewCList(
		&ParamsCList[string]{
			Processor: stringProcessor,
			Workers:   1,
		},
	)
	require.NoError(b, errCr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		list.Insert(fmt.Sprintf("item-%d", i))
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkProcess/workers-1-items-100-16         	    3397	    345445 ns/op	  138469 B/op	    1378 allocs/op
// BenchmarkProcess/workers-1-items-1000-16        	     214	   5665258 ns/op	10911919 B/op	   16868 allocs/op
// BenchmarkProcess/workers-1-items-10000-16       	       2	 737495478 ns/op	1560637828 B/op	  229071 allocs/op
// BenchmarkProcess/workers-2-items-100-16         	    4952	    227202 ns/op	   75812 B/op	     942 allocs/op
// BenchmarkProcess/workers-2-items-1000-16        	     319	   3711189 ns/op	 5529228 B/op	   10949 allocs/op
// BenchmarkProcess/workers-2-items-10000-16       	       3	 378127194 ns/op	781236925 B/op	  141281 allocs/op
// BenchmarkProcess/workers-4-items-100-16         	    5730	    187627 ns/op	   44488 B/op	     724 allocs/op
// BenchmarkProcess/workers-4-items-1000-16        	     477	   2500815 ns/op	 2836659 B/op	    7980 allocs/op
// BenchmarkProcess/workers-4-items-10000-16       	       5	 210088730 ns/op	391471875 B/op	   96760 allocs/op
// BenchmarkProcess/workers-8-items-100-16         	    8007	    151087 ns/op	   29880 B/op	     619 allocs/op
// BenchmarkProcess/workers-8-items-1000-16        	     685	   1725423 ns/op	 1490377 B/op	    6495 allocs/op
// BenchmarkProcess/workers-8-items-10000-16       	       9	 120080285 ns/op	196526385 B/op	   73972 allocs/op
// BenchmarkProcess/workers-16-items-100-16        	    9040	    123524 ns/op	   22514 B/op	     566 allocs/op
// BenchmarkProcess/workers-16-items-1000-16       	     933	   1256831 ns/op	  810516 B/op	    5755 allocs/op
// BenchmarkProcess/workers-16-items-10000-16      	      18	  66234018 ns/op	99005428 B/op	   62195 allocs/op
func BenchmarkProcess(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8, 16}
	itemCounts := []int{100, 1000, 10000}

	for _, workers := range workerCounts {
		for _, items := range itemCounts {
			b.Run(
				fmt.Sprintf("workers-%d-items-%d", workers, items),
				func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						b.StopTimer()

						list, _ := NewCList(
							&ParamsCList[string]{
								Processor: stringProcessor,
								Workers:   workers,
							},
						)

						// Populate list
						for j := 0; j < items; j++ {
							list.Insert(fmt.Sprintf("item-%d", j))
						}

						b.StartTimer()
					}
				},
			)
		}
	}
}
