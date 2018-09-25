package queue_test

import (
	"fmt"

	"github.com/cloud-spin/queue"
)

func Example() {
	q := queue.New()

	for i := 0; i < 5; i++ {
		q.Put(i)
	}

	for !q.IsEmpty() {
		v, _ := q.Get()
		fmt.Print(v)
	}

	// Output: 01234
}
