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
	DefaultInternalArraySize = 100
)

// Queue imeplements a dybamic FIFO queue. The retrieved elements are retrieved in the same order they were added.
// Queue can be safety used from concurrent Go routines.
type Queue interface {
	Peek() interface{}
	Get() interface{}
	Put(v interface{}) error
	IsEmpty() bool
}

// QueueImpl implements a Queue.
type QueueImpl struct {
	head  *Node
	tail  *Node
	mutex *sync.Mutex
}

// Node represents a Queue node.
type Node struct {
	i byte
	v []interface{}
	n *Node
}

// NewQueue initializes a new instance of Queue.
// Lock controls thread-safety. If lock is true, the queue implementation will be thread-safe; false otherwise.
func NewQueue() Queue {
	q := &QueueImpl{
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

	q.checkLock()
	if q.tail == nil {
		q.tail = NewNode()
		q.head = q.tail
	}
	if len(q.tail.v) >= DefaultInternalArraySize {
		q.tail.n = NewNode()
		q.tail = q.tail.n
	}

	q.tail.v = append(q.tail.v, v)
	q.checkUnlock()

	return nil
}

// Get retrieves and removes the next element from the queue.
// If the queue is empty, nil will be returned.
func (q *QueueImpl) Get() interface{} {
	q.checkLock()
	if q.isEmpty() {
		q.checkUnlock()
		return nil
	}

	v := q.head.v[q.head.i]
	q.head.i++
	if q.head.i >= DefaultInternalArraySize {
		q.head = q.head.n
	}
	q.checkUnlock()

	return v
}

// Peek retrieves the next element from the queue, but does not remove it from the queue.
func (q *QueueImpl) Peek() interface{} {
	q.checkLock()
	if q.isEmpty() {
		q.checkUnlock()
		return nil
	}

	v := q.head.v[q.head.i]
	q.checkUnlock()
	return v
}

// IsEmpty returns true if the queue is empty; false otherwise.
func (q *QueueImpl) IsEmpty() bool {
	q.checkLock()
	res := q.isEmpty()
	q.checkUnlock()
	return res
}

// IsEmpty returns true if the queue is empty; false otherwise.
func (q *QueueImpl) isEmpty() bool {
	return q.head == nil || q.head.i >= byte(len(q.head.v))
}

func (q *QueueImpl) checkLock() {
	q.mutex.Lock()
}

func (q *QueueImpl) checkUnlock() {
	q.mutex.Unlock()
}
