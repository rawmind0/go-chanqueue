package chanqueue

import (
)

// Worker represents the worker that executes the job
type Worker struct {
	Id 			int
	// A dispatcher pool of workers channels to register
	WorkerPool  chan chan Job
	// Worker channel to receive job requests
	WorkerChannel  chan Job
	// Quit channel
	Quit    	chan bool
}

func NewWorker(i int, workerPool chan chan Job) *Worker {
	return &Worker{
		Id: i,
		WorkerPool: workerPool,
		WorkerChannel: make(chan Job),
		Quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func(w Worker) {
		for {
			// register worker as idle
			w.Idle()
			select {
			// Job received
			case job := <-w.WorkerChannel:
				job.Start()
			// Stop signal received
			case <-w.Quit:
				return		
			}
		}
	}(w)
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	w.Quit <- true
}

func (w Worker) Idle() {
	w.WorkerPool <- w.WorkerChannel
}