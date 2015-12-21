package goncurrence

// ProcessHandler is an interface to execute the job process
// Worker executes Exec() error function
type ProcessHandler interface {
	Exec() error
}

// WorkerManager is workers manager
type WorkerManager struct {
	workerNum   int
	isUnorderd  bool
	addCount    int
	endCount    int
	errs        []error
	startSignal chan int
	endSignal   chan int
	process     chan ProcessHandler
	err         chan error
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
		workerNum:   workerNum,
		startSignal: make(chan int),
		endSignal:   make(chan int),
		process:     make(chan ProcessHandler),
		err:         make(chan error),
	}

	// create workers
	for i := 0; i < workerNum; i++ {
		go w.startWorker()
	}

	return w
}

func (w *WorkerManager) startWorker() {
	<-w.startSignal
	for p := range w.process {
		err := p.Exec()
		if err != nil {
			w.err <- err
		}
		w.endSignal <- 0
	}
}

// IsUnordered set true to isUnorderd flag
func (w *WorkerManager) IsUnordered() {
	w.isUnorderd = true
}

// Add adds a new process handler to be executed
func (w *WorkerManager) Add(p ProcessHandler) *WorkerManager {
	go func() {
		w.process <- p
	}()
	w.addCount++
	return w
}

// Run execute all processes
// If isUnorderd is true and occurred error, last stacked is returned
func (w *WorkerManager) Run() error {
	// Send start channel
	for i := 0; i < w.workerNum; i++ {
		w.startSignal <- 0
	}

	for {
		select {
		case err := <-w.err:
			w.errs = append(w.errs, err)
			if !w.isUnorderd {
				return err
			}
		case <-w.endSignal:
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
