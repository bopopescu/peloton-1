package async

import (
	"context"
	"sync"
)

const (
	// DefaultMaxWorkers of a Pool. See Pool.SetMaxWorkers for more info.
	DefaultMaxWorkers = 4
)

// PoolOptions for constructing a new Pool.
type PoolOptions struct {
	MaxWorkers int
}

// Pool structure for running up to a maximum number of jobs concurrently.
// The pool as an internal queue, such that all jobs added will be accepted
// but not run until it reached the front of the queue and a worker is free.
type Pool struct {
	options    PoolOptions
	lock       sync.Mutex
	queue      *Queue
	numWorkers int
	jobs       sync.WaitGroup
}

// NewPool returns a new pool, provided the PoolOptions.
func NewPool(o PoolOptions) *Pool {
	if o.MaxWorkers <= 0 {
		o.MaxWorkers = DefaultMaxWorkers
	}

	p := &Pool{
		options:    o,
		queue:      NewQueue(),
		numWorkers: o.MaxWorkers,
	}

	// Spawn initial workers.
	for i := 0; i < o.MaxWorkers; i++ {
		go p.runWorker()
	}

	return p
}

// SetMaxWorkers to the number provided. If smaller than the current value, it
// will lazily close existing workers. If greater, new workers will be created.
// If 0 or less is given, DefaultMaxWorkers will be used instead.
func (p *Pool) SetMaxWorkers(num int) {
	if num <= 0 {
		num = DefaultMaxWorkers
	}

	p.lock.Lock()
	oldNum := p.numWorkers
	p.numWorkers = num
	p.lock.Unlock()

	for i := oldNum; i < num; i++ {
		go p.runWorker()
	}
}

// Enqueue a job in the pool.
// TODO: Take an context argument that will be associated to the job. That way
// deadlines can easily be propagated.
func (p *Pool) Enqueue(job Job) {
	p.jobs.Add(1)
	p.queue.Enqueue(job)
}

// WaitUntilProcessed will block until both the queue is empty and all workers
// are idle. This is useful for per-request Pools and in testing.
func (p *Pool) WaitUntilProcessed() {
	p.jobs.Wait()
}

func (p *Pool) runWorker() {
	for {
		if p.shouldWorkerStop() {
			return
		}

		job := <-p.queue.DequeueChannel()
		// TODO: Implement Stop() on queue that allows termination of all jobs.
		job.Run(context.TODO())
		p.jobs.Done()
	}
}

func (p *Pool) shouldWorkerStop() bool {
	stop := false
	p.lock.Lock()
	if p.numWorkers > p.options.MaxWorkers {
		p.numWorkers--
		stop = true
	}
	p.lock.Unlock()
	return stop
}
