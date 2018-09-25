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
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewQueuShouldReturnInitiazedInstanceOfQueue(t *testing.T) {
	q := New()

	if q == nil {
		t.Error("Expected: new instance of queue; Got: nil")
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
			putCount:       []int{100},
			getCount:       []int{100},
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
			q := New()
			lastPut := 0
			lastGet := 0
			for count := 0; count < len(test.getCount); count++ {
				for i := 1; i <= test.putCount[count]; i++ {
					lastPut++
					q.Put(lastPut)
					if v, ok := q.Peek(); !ok || v != lastGet+1 {
						t.Errorf("Expected: %d; Got: %d", lastGet, v)
					}
				}
				for i := 1; i <= test.getCount[count]; i++ {
					lastGet++
					v, ok := q.Peek()
					if !ok || v.(int) != lastGet {
						t.Errorf("Expected: %d; Got: %d", lastGet, v)
					}
					v, ok = q.Get()
					if !ok || v.(int) != lastGet {
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

				v, ok := q.Peek()
				if !ok || v.(int) != lastGet {
					t.Errorf("Expected: %d; Got: %d", lastGet, v)
				}
				v, ok = q.Get()
				if !ok || v.(int) != lastGet {
					t.Errorf("Expected: %d; Got: %d", lastGet, v)
				}
			}
			v, ok := q.Peek()
			if ok || v != nil {
				t.Errorf("Expected: nil as the queue should be empty; Got: %d", v)
			}
			v, ok = q.Get()
			if ok || v != nil {
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
	q := New()
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
				v, ok := q.Get()
				if ok && v != nil {
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
