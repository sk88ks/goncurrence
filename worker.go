package goncurrency

// ProcessFunc is an function that would be executed as a process
type ProcessFunc func() error

// WorkerManager is workers manager
type WorkerManager struct {
	workerNum  int
	isUnorderd bool
	addCount   int
	endCount   int
	errs       []error
	do         chan struct{}
	done       chan struct{}
	end        chan struct{}
	process    chan ProcessFunc
	err        chan error
}

// DefaultProcess is default process implementing process handler
type DefaultProcess struct {
	Func   func() (interface{}, error)
	Result interface{}
}

// Exec executes job process
func (d *DefaultProcess) Exec() error {
	res, err := d.Func()
	if err != nil {
		return err
	}

	d.Result = res

	return nil
}

// New creates a new worker manager
func New(workerNum int) *WorkerManager {
	w := &WorkerManager{
		workerNum: workerNum,
		do:        make(chan struct{}),
		done:      make(chan struct{}),
		end:       make(chan struct{}),
		process:   make(chan ProcessFunc),
		err:       make(chan error),
	}

	// create workers
	for i := 0; i < workerNum; i++ {
		go w.startWorker()
	}

	return w
}

func (w *WorkerManager) startWorker() {
	<-w.do
	for {
		select {
		case p := <-w.process:
			err := p()
			if err != nil {
				w.err <- err
			}
			w.end <- struct{}{}
		case <-w.done:
			return
		}
	}
}

// IsUnordered set true to isUnorderd flag
func (w *WorkerManager) IsUnordered() {
	w.isUnorderd = true
}

// Add adds a new process handler to be executed
func (w *WorkerManager) Add(ps ...ProcessFunc) *WorkerManager {
	for i := range ps {
		go func() {
			w.process <- ps[i]
		}()
		w.addCount++
	}
	return w
}

// Run execute all processes
// If isUnorderd is true and occurred error, last stacked is returned
func (w *WorkerManager) Run() error {
	// Start all processes
	close(w.do)

	for {
		select {
		case err := <-w.err:
			w.errs = append(w.errs, err)
			if !w.isUnorderd {
				return err
			}
		case <-w.end:
			w.endCount++
		}

		if w.endCount == w.addCount {
			break
		}
	}

	l := len(w.errs)
	if l > 0 {
		return w.errs[l-1]
	}

	return nil
}

// Errs returns all stacked errors
func (w *WorkerManager) Errs() []error {
	return w.errs
}
