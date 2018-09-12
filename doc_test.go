package queue_test

import (
	"fmt"

	"github.com/cloud-spin/queue"
)

func Example() {
	q := queue.NewQueue()

	for i := 0; i < 5; i++ {
		q.Put(i)
	}

	for !q.IsEmpty() {
		fmt.Print(q.Get())
	}

	// Output: 01234
}
