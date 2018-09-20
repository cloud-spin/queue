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
	"errors"
	"sync"
)

const (
	// DefaultInternalArraySize holds the size of each internal array.
	DefaultInternalArraySize = 128
)

// Queue represents a thread-safe, dynamically growing FIFO queue.
type Queue interface {
	Get() interface{}
	Put(v interface{}) error
	Peek() interface{}
	IsEmpty() bool
}

// QueueImpl implements a Queue.
type QueueImpl struct {
	head  *Node
	tail  *Node
	pos   int
	mutex *sync.Mutex
}

// Node represents a Queue node.
type Node struct {
	v []interface{}
	n *Node
}

// NewQueue initializes a new instance of Queue.
func NewQueue() Queue {
	head := NewNode()
	q := &QueueImpl{
		head:  head,
		tail:  head,
		mutex: &sync.Mutex{},
	}
	return q
}

// NewNode initializes a new instance of Node.
func NewNode() *Node {
	return &Node{
		v: make([]interface{}, 0, DefaultInternalArraySize),
	}
}

// Put adds a value to the queue.
// Put doesn't accept nil values.
func (q *QueueImpl) Put(v interface{}) error {
	if v == nil {
		return errors.New("Cannot add nil value")
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	if len(q.tail.v) >= DefaultInternalArraySize {
		q.tail.n = NewNode()
		q.tail = q.tail.n
	}
	q.tail.v = append(q.tail.v, v)

	return nil
}

// Get retrieves and removes the next element from the queue.
// If the queue is empty, nil will be returned.
func (q *QueueImpl) Get() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.isEmpty() {
		return nil
	}

	v := q.head.v[q.pos]
	q.head.v[q.pos] = nil
	q.pos++
	if q.pos >= DefaultInternalArraySize {
		q.head = q.head.n
		q.pos = 0
	}

	return v
}

// Peek retrieves the next element from the queue, but does not remove it from the queue.
func (q *QueueImpl) Peek() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.isEmpty() {
		return nil
	}

	v := q.head.v[q.pos]
	return v
}

// IsEmpty returns true if the queue is empty; false otherwise.
func (q *QueueImpl) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	res := q.isEmpty()
	return res
}

func (q *QueueImpl) isEmpty() bool {
	return q.head == nil || q.pos >= len(q.head.v)
}
