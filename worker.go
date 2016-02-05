package goncurrency

import "reflect"

// ProcessFunc is an function that would be executed as a process
type ProcessFunc func() (interface{}, error)

// WorkerManager is workers manager
type WorkerManager struct {
	workerNum int
	addCount  int
	done      chan struct{}
	process   chan ProcessFunc
	result    chan ProcessResult
}

// ProcessResult is result of processes
type ProcessResult struct {
	v   reflect.Value
	err error
}

// ProcessIterator is iterator of processes results
type ProcessIterator struct {
	wm       *WorkerManager
	endCount int
	result   *ProcessResult
}

// New creates a new worker manager
func New(workerNum int) *WorkerManager {
	w := &WorkerManager{
		workerNum: workerNum,
		done:      make(chan struct{}),
		process:   make(chan ProcessFunc),
		result:    make(chan ProcessResult),
	}

	// create workers
	for i := 0; i < workerNum; i++ {
		go w.startWorker()
	}

	return w
}

func (w *WorkerManager) startWorker() {
	for {
		select {
		case p := <-w.process:
			res, err := p()
			w.result <- ProcessResult{
				v:   reflect.ValueOf(res),
				err: err,
			}
		case <-w.done:
			return
		}
	}
}

// Add adds a new process handler to be executed
func (w *WorkerManager) Add(ps ...ProcessFunc) *WorkerManager {
	for i := range ps {
		go func(pf ProcessFunc) {
			w.process <- pf
		}(ps[i])
		w.addCount++
	}
	return w
}

// Iter gets iterator for processes results
func (w *WorkerManager) Iter() *ProcessIterator {
	return &ProcessIterator{
		wm: w,
	}
}

// Next iterate processes results
func (iter *ProcessIterator) Next() bool {
	if iter.wm.addCount <= iter.endCount {
		close(iter.wm.done)
		close(iter.wm.result)
		return false
	}
	res := <-iter.wm.result
	iter.result = &res
	iter.endCount++
	return true
}

// Result set result to destination data
func (iter *ProcessIterator) Result(dst interface{}) error {
	if iter.result.err != nil {
		return iter.result.err
	}

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return nil
	}

	dstValue.Elem().Set(iter.result.v)
	return nil
}
