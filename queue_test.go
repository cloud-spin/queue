// Copyright (c) 2018 cloud-spin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package queue

import (
	"container/list"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewQueuShouldReturnInitiazedInstanceOfQueue(t *testing.T) {
	q := NewQueue()

	if q == nil {
		t.Error("Expected: new instance of queue; Got: nil")
	}
}

func TestPutNilValueShouldReturnError(t *testing.T) {
	q := NewQueue()
	if err := q.Put(nil); err == nil {
		t.Error("Expected: error; Got: success")
	}
}

func TestPutGetAndPeekShouldRetrieveAllElementsInOrder(t *testing.T) {
	tests := map[string]struct {
		putCount       []int
		getCount       []int
		remainingCount int
	}{
		"Test 1 item": {
			putCount:       []int{1},
			getCount:       []int{1},
			remainingCount: 0,
		},
		"Test 100 items": {
			putCount:       []int{2},
			getCount:       []int{2},
			remainingCount: 0,
		},
		"Test 1000 items": {
			putCount:       []int{1000},
			getCount:       []int{1000},
			remainingCount: 0,
		},
		"Test sequence 1": {
			putCount:       []int{1, 2, 100, 101},
			getCount:       []int{1, 2, 100, 101},
			remainingCount: 0,
		},
		"Test sequence 2": {
			putCount:       []int{10, 1},
			getCount:       []int{1, 10},
			remainingCount: 0,
		},
		"Test sequence 3": {
			putCount:       []int{101, 101},
			getCount:       []int{100, 101},
			remainingCount: 1,
		},
		"Test sequence 4": {
			putCount:       []int{1000, 1000, 1001},
			getCount:       []int{10, 10, 1},
			remainingCount: 2980,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("running tests %s", name)
			q := NewQueue()
			lastPut := 0
			lastGet := 0
			for count := 0; count < len(test.getCount); count++ {
				for i := 1; i <= test.putCount[count]; i++ {
					lastPut++
					if err := q.Put(lastPut); err != nil {
						t.Errorf("Expected: no error; Got: %s", err.Error())
					}
					if v := q.Peek(); v != lastGet+1 {
						t.Errorf("Expected: %d; Got: %d", lastGet, v)
					}
				}
				for i := 1; i <= test.getCount[count]; i++ {
					lastGet++
					v := q.Peek().(int)
					if v != lastGet {
						t.Errorf("Expected: %d; Got: %d", lastGet, v)
					}
					v = q.Get().(int)
					if v != lastGet {
						t.Errorf("Expected: %d; Got: %d", lastGet, v)
					}
				}
			}

			if test.remainingCount > 0 {
				if q.IsEmpty() {
					t.Error("Expected: non-empty queue; Got: empty")
				}
			} else {
				if !q.IsEmpty() {
					t.Error("Expected: empty queue; Got: non-empty")
				}
			}

			for i := 1; i <= test.remainingCount; i++ {
				lastGet++

				v := q.Peek().(int)
				if v != lastGet {
					t.Errorf("Expected: %d; Got: %d", lastGet, v)
				}
				v = q.Get().(int)
				if v != lastGet {
					t.Errorf("Expected: %d; Got: %d", lastGet, v)
				}
			}
			v := q.Peek()
			if v != nil {
				t.Errorf("Expected: nil as the queue should be empty; Got: %d", v)
			}
			v = q.Get()
			if v != nil {
				t.Errorf("Expected: nil as the queue should be empty; Got: %d", v)
			}

			if !q.IsEmpty() {
				t.Error("Expected: empty queue; Got: non-empty")
			}
		})
	}
}

func TestPutAndGetConcurrently(t *testing.T) {
	const routines = 5
	q := NewQueue()
	values := map[int]int{}
	wgPut := sync.WaitGroup{}
	wgPut.Add(routines)
	wgGet := sync.WaitGroup{}
	wgGet.Add(routines)
	var waitGet int32
	mux := &sync.Mutex{}

	for i := 0; i < routines; i++ {
		go func() {
			for v := 1; v < 1000; v++ {
				q.Put(v)
			}
			wgPut.Done()
		}()
	}

	for i := 0; i < routines; i++ {
		go func() {
			for {
				v := q.Get()
				if v != nil {
					iv := v.(int)
					mux.Lock()
					if count, ok := values[iv]; ok {
						values[iv] = count + 1
					} else {
						values[iv] = 1
					}
					mux.Unlock()
				}

				if atomic.LoadInt32(&waitGet) != 0 && q.IsEmpty() {
					wgGet.Done()
					break
				}
			}
		}()
	}

	wgPut.Wait()
	atomic.StoreInt32(&waitGet, 1)
	wgGet.Wait()

	for _, v := range values {
		if v != routines {
			t.Errorf("Expected: %d as there were %d routines adding the same numbers to the queue; Got: %d", routines, routines, v)
		}
	}
}

func BenchmarkPutAndGet(b *testing.B) {
	q := NewQueue()

	for n := 0; n < b.N; n++ {
		q.Put(n)
	}

	lastGet := 0
	for !q.IsEmpty() {
		v := q.Get()
		if v == nil {
			b.Errorf("Expected: %d; Got: nil", lastGet)
		} else if v.(int) != lastGet {
			b.Errorf("Expected: %d; Got: %d", lastGet, v.(int))
		}
		lastGet++
	}
}

func BenchmarkStandardListPackage(b *testing.B) {
	l := list.New()

	for n := 0; n < b.N; n++ {
		l.PushBack(n)
	}

	lastGet := 0
	for e := l.Front(); e != nil; e = e.Next() {
		v := e.Value
		if v == nil {
			b.Errorf("Expected: %d; Got: nil", v)
		} else if v.(int) != lastGet {
			b.Errorf("Expected: %d; Got: %d", lastGet, v.(int))
		}
		lastGet++
	}
}
