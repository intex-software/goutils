package queue

import (
	"sync"
)

type ConcurrentQueue[T any] struct {
	items []T

	lock *sync.Mutex
	cond *sync.Cond
}

func NewConcurrent[T any]() *ConcurrentQueue[T] {
	lock := &sync.Mutex{}

	return &ConcurrentQueue[T]{
		lock: lock,
		cond: sync.NewCond(lock),
	}
}

func (q *ConcurrentQueue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *ConcurrentQueue[T]) Dequeue() (item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait()
	}

	item, q.items = q.items[0], q.items[1:]
	return
}

func (q *ConcurrentQueue[T]) Empty() bool {
	return len(q.items) == 0
}
