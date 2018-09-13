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
	"testing"
)

func BenchmarkQueuePackage(b *testing.B) {
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
