# Queue [![Build Status](https://travis-ci.com/cloud-spin/queue.svg?branch=master)](https://travis-ci.com/cloud-spin/queue) [![codecov](https://codecov.io/gh/cloud-spin/queue/branch/master/graph/badge.svg)](https://codecov.io/gh/cloud-spin/queue) [![Go Report Card](https://goreportcard.com/badge/github.com/cloud-spin/queue)](https://goreportcard.com/report/github.com/cloud-spin/queue)  [![GoDoc](https://godoc.org/github.com/cloud-spin/queue?status.svg)](https://godoc.org/github.com/cloud-spin/queue)

Package queue implements a safe for concurrent use and dynamically growing queue that uses linked arrays as its internal data structure.

## Install
From a configured [Go environment](https://golang.org/doc/install#testing):
```sh
go get -u github.com/cloud-spin/queue
```

If you are using dep:
```sh
dep ensure -add github.com/cloud-spin/queue
```

## How to Use
```go
package main

import (
	"fmt"

	"github.com/cloud-spin/queue"
)

func main() {
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


## Performance
Package queue implements a FIFO queue storing the elements in linked arrays.

Why?

Locality is the answer! Traditional queue implementations uses linked list as its internal data structure.
While linked lists are great for dynamically growing lists with little overhead, navigating and retrieving the elements
is tipically slower than retrieving elements from arrays because linked lists require "jumps" (next element) to potentially far away memory addresses. Arrays, on the other hand, doesn't suffer from the same problem as the elements are stored in sequential memory addresses, making it faster to navigate and retrieve the subsequent elements. By using a linked list, linking fixed sized arrays together, package queue is able to take advantage of the memory cache locality of the arrays and is still able to dynamically grow with little overhead.

The below [benchmark test](benchmark_test.go) result adds 100 items to the queue and then removes them afterwards.

```
go version
go version go1.11 darwin/amd64

sysctl -n machdep.cpu.brand_string
Intel(R) Core(TM) i5-7267U CPU @ 3.10GHz

go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/cloud-spin/queue
BenchmarkQueuePackage-4   	  100000	     17002 ns/op	    2912 B/op	     103 allocs/op
PASS
ok  	github.com/cloud-spin/queue	1.884s
```

These tests can be found at [benchmark tests](benchmark_test.go).



## Releases
We're committed to a CI/CD lifecycle releasing frequent, but only stable, production ready versions with all proper tests in place.

We strive as much as possible to keep backwards compatibility with previous versions, so breaking changes are a no-go.

For a list of changes in each released version, see [CHANGELOG.md](CHANGELOG.md).


## Supported Go Versions
See [supported_go_versions.md](https://github.com/cloud-spin/docs/blob/master/supported_go_versions.md).


## License
MIT, see [LICENSE](LICENSE).

"Use, abuse, have fun and contribute back!"


## Contributions
See [contributing.md](https://github.com/cloud-spin/docs/blob/master/contributing.md).
