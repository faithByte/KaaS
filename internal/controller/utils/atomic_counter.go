package utils

import "sync/atomic"

type AtomicCounter struct {
	counter int32
}

func (data *AtomicCounter) Increment() {
	atomic.AddInt32(&data.counter, 1)
}

func (data *AtomicCounter) Decrement() {
	atomic.AddInt32(&data.counter, -1)
}

func (data *AtomicCounter) Value() int32 {
	return atomic.LoadInt32(&data.counter)
}
