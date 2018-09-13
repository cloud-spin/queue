# Queue [![Build Status](https://travis-ci.com/cloud-spin/queue.svg?branch=master)](https://travis-ci.com/cloud-spin/queue) [![codecov](https://codecov.io/gh/cloud-spin/queue/branch/master/graph/badge.svg)](https://codecov.io/gh/cloud-spin/queue) [![Go Report Card](https://goreportcard.com/badge/github.com/cloud-spin/queue)](https://goreportcard.com/report/github.com/cloud-spin/queue)  [![GoDoc](https://godoc.org/github.com/cloud-spin/queue?status.svg)](https://godoc.org/github.com/cloud-spin/queue)

Package queue implements a high performance, thread-safe and dynamically growing queue that uses linked arrays as its internal data structure.
Package queue is about 30% faster than the standard [list package](https://github.com/golang/go/tree/master/src/container/list) and uses 50% less memory, besides being thread-safe.*

#### Install

From a configured [Go environment](https://golang.org/doc/install#testing):
```sh
go get -u github.com/cloud-spin/queue
```

If you are using dep:
```sh
dep ensure -add github.com/cloud-spin/queue
```


#### How to Use

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

#### Performance
Package queue implements a FIFO queue storing the elements in linked arrays.

Why?

Locality is the answer! Traditional queue implementations uses linked list as its internal data structure.
While linked lists are great for dynamically growing lists with little overhead, navigating and retrieving the elements
is tipically slower than retrieving elements from arrays because linked lists require "jumps" (next element) to potentially far away memory addresses. Arrays, on the other hand, doesn't suffer from the same problem as the elements are stored in sequential memory addresses, making it faster to navigate and retrieve the subsequent elements. By using a linked list, linking fixed sized arrays together, package queue is able to take advantage of the memory cache locality of the arrays and is still able to dynamically grow with little overhead.

As Go doesn't provide a standard queue implementation, the likely closest standand Go structure to a queue is the [list package](https://github.com/golang/go/tree/master/src/container/list). Below benchmark tests compare adding and retrieving elements with the list package.

```
go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/cloud-spin/queue
BenchmarkQueuePackage-4          	20000000	        95.9 ns/op	      26 B/op	       1 allocs/op
BenchmarkStandardListPackage-4   	10000000	       135 ns/op	      56 B/op	       2 allocs/op
PASS
ok  	github.com/cloud-spin/queue	3.607s
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


## License
MIT, see [LICENSE](LICENSE).

"Use, abuse, have fun and contribute back!"



## Contributions
By reading, analysing and using many of the tons of super high quality and amazing open source projects out there, we feel we learned so much from the wonderful open source community that we strongly feel the need to contribute back!

CloudSpin is our way to say "Thank you very much open source community for all the support for software development!".

As a member of the open source community, please help grow the community by using, sharing and contributing back to any of the tons of open source projects out there. The community altogether is super excited to see new, engaged members, and so are we at CloudSpin!

If you are interested in contributing to CloudSpin, for now help to spread the word and let others know about CloudSpin. Once CloudSpin is fully tested and verified to be working as expected, we'll start accepting contributions.

We're 100% commited to below software development rules for CloudSpin:

- Simplicity
- Testable code
- Tests, tests, tests!
	- Strong test suite covering all major code routes/branches
	- Strong focus to achieve 100% code coverage everywhere

On top of that, clean code is a must. No weird, obscure logic anywhere. As part of that principle, we avoid using comments to describe code logic as we strive to make the code so clean that any inline comments would just pollute the code (no need to explain what is very clear already!). The rule we follow is this: if we feel the need to add a comment to explain something, think again. That likely means the logic is too complex. Rethink and simplify it!



*According to the benchmark tests in this package.
