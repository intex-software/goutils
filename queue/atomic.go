package queue

import (
	"sync/atomic"
)

const (
	increment = 1
	decrement = ^uint64(0)
)

type AtomicQueue[T any] struct {
	items    chan *T
	counter  *uint64
	capacity uint64
}

func NewAtomic[T any](capacity uint64) *AtomicQueue[T] {
	return &AtomicQueue[T]{
		items:    make(chan *T, capacity),
		counter:  new(uint64),
		capacity: capacity,
	}
}

func (q *AtomicQueue[T]) Enqueue(item *T) {
	atomic.AddUint64(q.counter, increment)
	q.items <- item
}

func (q *AtomicQueue[T]) Dequeue() *T {
	item := <-q.items
	atomic.AddUint64(q.counter, decrement)
	return item
}

func (q *AtomicQueue[T]) Empty() bool {
	return *q.counter == 0
}

func (q *AtomicQueue[T]) Full() bool {
	return *q.counter == uint64(q.capacity)
}

func (q *AtomicQueue[T]) Capacity() uint64 {
	return q.capacity
}

func (q *AtomicQueue[T]) Size() int {
	return int(*q.counter)
}
