package render

import "sync"

type Future[T any] struct {
	mu     sync.Mutex
	ready  chan struct{}
	result T
	err    error
	set    bool
}

func NewFuture[T any]() *Future[T] {
	return &Future[T]{ready: make(chan struct{})}
}

func (f *Future[T]) Set(result T, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.set {
		return
	}

	f.result = result
	f.err = err
	f.set = true
	close(f.ready)
}

func (f *Future[T]) Get() (T, error) {
	<-f.ready
	return f.result, f.err
}

func (f *Future[T]) IsDone() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.set
}
