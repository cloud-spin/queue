# Queue [![Build Status](https://travis-ci.com/cloud-spin/queue.svg?branch=master)](https://travis-ci.com/cloud-spin/queue) [![codecov](https://codecov.io/gh/cloud-spin/queue/branch/master/graph/badge.svg)](https://codecov.io/gh/cloud-spin/queue) [![Go Report Card](https://goreportcard.com/badge/github.com/cloud-spin/queue)](https://goreportcard.com/report/github.com/cloud-spin/queue)  [![GoDoc](https://godoc.org/github.com/cloud-spin/queue?status.svg)](https://godoc.org/github.com/cloud-spin/queue)

Package queue implements a high performance, thread-safe and dynamically growing queue that uses linked arrays as its internal data structure.
Package queue is about 30% faster than the standard [list package](https://github.com/golang/go/tree/master/src/container/list) and uses 50% less memory.

#### How to Use

```go
package main

import (
	"fmt"

	"github.com/cloud-spin/queue"
)

func Example() {
	q := queue.NewQueue()

	for i := 1; i <= 5; i++ {
		q.Put(i)
	}

	for !q.IsEmpty() {
		fmt.Println(q.Get())
	}
}
```

Output:
```
1
2
3
4
5
```

Also refer to the tests at [queue_test.go](queue_test.go).

#### Performance
Package queue implements a FIFO queue storing the elements in linked arrays.

Why?

Locality is the answer! Traditional queue implementations uses linked list as its internal data structure.
While linked lists are great for dynamically growing lists with little overhead, navigating and retrieving the elements
is tipically slower than retrieving elements from arrays because linked lists require "jumps" from potentially very distinct and unkown
memory addresses. Arrays, on the other hand, doesn't suffer from the same problem as the elements are stored in sequential memory addresses,
making it faster to navigate and retrieve the subsequent elements.

As Go doesn't provide a standard Queue implementation, the likely closest standand Go structure to a Queue is the [list package](https://github.com/golang/go/tree/master/src/container/list). Below benchmark tests compare adding and retrieving elements with the list package.

```
go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/cloud-spin/queue
BenchmarkPutAndGet-4             	20000000	        94.6 ns/op	      26 B/op	       1 allocs/op
BenchmarkStandardListPackage-4   	20000000	       139 ns/op	      56 B/op	       2 allocs/op
```

These tests can be found at [queue_test.go](queue_test.go).

### How To Run Tests
Run below commands from the package root directory.

Run tests with code coverage:
```
go test -coverprofile=coverage.txt -covermode=atomic
```

Get code coverage report:
```
go tool cover -html=coverage.txt
```

Run bench tests:
```
go test -bench=. -benchmem
```
