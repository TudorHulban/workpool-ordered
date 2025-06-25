package workpoolordered

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Simple string processor for benchmarking
func stringProcessor(s string) (string, error) {
	// Simulate some work
	return fmt.Sprintf("processed-%s", s), nil
}

// Heavy computation processor
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
	list := NewDLinkedList(stringProcessor, 1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		list.Insert(fmt.Sprintf("item-%d", i))
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkProcess/workers-1-items-100-16         	   42091	     27662 ns/op	    7899 B/op	     312 allocs/op
// BenchmarkProcess/workers-1-items-1000-16        	    4161	    274117 ns/op	   73733 B/op	    3015 allocs/op
// BenchmarkProcess/workers-1-items-10000-16       	     453	   2624116 ns/op	  870773 B/op	   30024 allocs/op
// BenchmarkProcess/workers-2-items-100-16         	   39850	     30560 ns/op	    7979 B/op	     313 allocs/op
// BenchmarkProcess/workers-2-items-1000-16        	    4518	    271812 ns/op	   73851 B/op	    3016 allocs/op
// BenchmarkProcess/workers-2-items-10000-16       	     460	   2559466 ns/op	  870855 B/op	   30025 allocs/op
// BenchmarkProcess/workers-4-items-100-16         	   36756	     32754 ns/op	    8135 B/op	     315 allocs/op
// BenchmarkProcess/workers-4-items-1000-16        	    4341	    281879 ns/op	   74047 B/op	    3019 allocs/op
// BenchmarkProcess/workers-4-items-10000-16       	     452	   2646399 ns/op	  871070 B/op	   30028 allocs/op
// BenchmarkProcess/workers-8-items-100-16         	   33566	     35715 ns/op	    8433 B/op	     319 allocs/op
// BenchmarkProcess/workers-8-items-1000-16        	    4497	    270092 ns/op	   74419 B/op	    3023 allocs/op
// BenchmarkProcess/workers-8-items-10000-16       	     441	   2815325 ns/op	  874101 B/op	   30047 allocs/op
// BenchmarkProcess/workers-16-items-100-16        	   30008	     40245 ns/op	    9029 B/op	     327 allocs/op
// BenchmarkProcess/workers-16-items-1000-16       	    4268	    274955 ns/op	   75089 B/op	    3032 allocs/op
// BenchmarkProcess/workers-16-items-10000-16      	     412	   2922570 ns/op	  875219 B/op	   30057 allocs/op
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

						list := NewDLinkedList(stringProcessor, workers)

						// Populate list
						for j := 0; j < items; j++ {
							list.Insert(fmt.Sprintf("item-%d", j))
						}

						b.StartTimer()
						list.Process(context.Background())
					}
				},
			)
		}
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkProcessHeavy/heavy-workers-1-16         	    1366	    867895 ns/op	   81874 B/op	    4015 allocs/op
// BenchmarkProcessHeavy/heavy-workers-2-16         	    1572	    761309 ns/op	   81978 B/op	    4016 allocs/op
// BenchmarkProcessHeavy/heavy-workers-4-16         	    2022	    613588 ns/op	   82177 B/op	    4019 allocs/op
// BenchmarkProcessHeavy/heavy-workers-8-16         	    2148	    561129 ns/op	   82582 B/op	    4024 allocs/op
func BenchmarkProcessHeavy(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8}
	items := 1000

	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("heavy-workers-%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				list := NewDLinkedList(heavyProcessor, workers)

				for j := 0; j < items; j++ {
					list.Insert(fmt.Sprintf("item-%d", j))
				}

				b.StartTimer()
				list.Process(context.Background())
			}
		})
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkProcessIO/io-workers-1-16         	      10	 105612461 ns/op	    7315 B/op	     313 allocs/op
// BenchmarkProcessIO/io-workers-10-16        	     100	  10703209 ns/op	    9337 B/op	     333 allocs/op
// BenchmarkProcessIO/io-workers-50-16        	     543	   2216573 ns/op	   15733 B/op	     412 allocs/op
// BenchmarkProcessIO/io-workers-100-16       	    1068	   1154229 ns/op	   24325 B/op	     513 allocs/op
func BenchmarkProcessIO(b *testing.B) {
	workerCounts := []int{1, 10, 50, 100}
	items := 100

	for _, workers := range workerCounts {
		b.Run(
			fmt.Sprintf("io-workers-%d", workers),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					list := NewDLinkedList(ioProcessor, workers)

					for j := 0; j < items; j++ {
						list.Insert(fmt.Sprintf("item-%d", j))
					}

					b.StartTimer()
					list.Process(context.Background())
				}
			},
		)
	}
}

func BenchmarkRead(b *testing.B) {
	itemCounts := []int{10, 100, 200}

	for _, items := range itemCounts {
		b.Run(
			fmt.Sprintf("items-%d", items),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					list := NewDLinkedList(stringProcessor, 4)

					// Populate and process
					for j := 0; j < items; j++ {
						list.Insert(fmt.Sprintf("item-%d", j))
					}
					list.Process(context.Background())

					b.StartTimer()
					list.Read()
				}
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkMemoryAllocation/current-implementation-16         	  942913	      1143 ns/op	     392 B/op	      12 allocs/op
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run(
		"current-implementation",
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				list := NewDLinkedList(stringProcessor, 4)

				for j := 0; j < 1; j++ {
					list.Insert(fmt.Sprintf("item-%d", j))
				}

				list.Process(context.Background())
				list.Read()
			}
		},
	)
}

// BenchmarkConcurrentAccess-16    	 3194024	       380.0 ns/op	     201 B/op	       6 allocs/op
func BenchmarkConcurrentAccess(b *testing.B) {
	list := NewDLinkedList(stringProcessor, 4)

	b.RunParallel(
		func(pb *testing.PB) {
			i := 0

			for pb.Next() {
				// Mix of operations
				switch i % 3 {
				case 0:
					list.Insert(fmt.Sprintf("item-%d", i))
				case 1:
					list.Process(context.Background())
				case 2:
					list.Read()
				}

				i++
			}
		},
	)
}
