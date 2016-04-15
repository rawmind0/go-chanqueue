package chanqueue

import (
	"fmt"
)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	ChanPool chan chan Job
	// Input queue. Dispatch received jobs from it
	InputQueue chan Job
	// slice of instantiated workers
	Workers []*Worker
	// MaxWorkers entries
	MaxWorkers int
	// Quit channel
	Quit chan bool
	Status string
}

func NewDispatcher(iQueue chan Job, maxWorkers int) *Dispatcher {
	return &Dispatcher{
		ChanPool: make(chan chan Job, maxWorkers), 
		InputQueue: iQueue, 
		MaxWorkers: maxWorkers, 
		Workers: make([]*Worker,maxWorkers),
		Quit: make(chan bool)}
}

func (d *Dispatcher) Start() {
	d.SetStatus("Starting")
    fmt.Printf("[dispatcher]: Starting...\n")
    d.PrefetchWorkers()
	go d.DispatchJobs()
}

// Listen jobs from input queue
func (d *Dispatcher) DispatchJobs() {
	d.SetStatus("Listening")
	fmt.Printf("[dispatcher]: Listening from input queue\n")
	for {
		select {
		case job := <- d.InputQueue:
			// a job request has been received thought InputQueue 
			go func(d *Dispatcher, job Job) {
				// try to obtain a worker channel from ChanPool and dispatch the job
				// this will block until a worker is idle
				WorkerChannel := <-d.ChanPool
				WorkerChannel <- job
			}(d, job)
		case <-d.Quit:
			// we have received a signal to stop
			return
		}
	}
}

// Stop signals the worker to stop listening for work requests.
func (d *Dispatcher)  Stop() {
	d.SetStatus("Stoping")
	fmt.Printf("[dispatcher]: Stoping dispatcher...\n")
	for _,worker := range d.Workers {
		worker.Stop()
	}
	go func(d *Dispatcher) {
		d.Quit <- true
	}(d)
}

// Prefecth MaxWorkers workers channels
func (d *Dispatcher) PrefetchWorkers() {
	d.SetStatus("Prefetching")
	fmt.Printf("[dispatcher]: Prefetching %d workers...\n", d.MaxWorkers)

	for i := 0; i < d.MaxWorkers; i++ {
		d.Workers[i] = NewWorker(i, d.ChanPool)
		d.Workers[i].Start()
	}
}

func (d *Dispatcher) SetStatus(s string) {
	d.Mutex.Lock()
	d.Status = s
	d.Mutex.Unlock()
}