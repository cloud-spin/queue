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
	"testing"
)

const (
	// testCount hold the number of items to add to the queue.
	testCount = 100
)

func BenchmarkQueuePackage(b *testing.B) {
	for n := 0; n < b.N; n++ {
		q := New()
		lastGet := 0

		for i := 0; i < testCount; i++ {
			q.Put(i)
		}
		for !q.IsEmpty() {
			v, ok := q.Get()
			if !ok || v == nil {
				b.Errorf("Expected: %d; Got: nil", lastGet)
			} else if v.(int) != lastGet {
				b.Errorf("Expected: %d; Got: %d", lastGet, v.(int))
			}
			lastGet++
		}
	}
}
