package workpoolordered

import (
	"fmt"
	"time"
)

func (dl *DLinkedList[T]) worker() {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
