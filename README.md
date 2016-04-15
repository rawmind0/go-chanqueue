# Go chanqueue package

This repository contains package to use system to enqueue and execute jobs asynchronously, 
using goroutines and channels for high concurrency and performance.

## How it Works

The queue system provides the abilty to enqueue and execute jobs (sync or async) with 
controlled concurrency in a very easy way. The package make use of goroutines and channels 
trying to provide best performance 

It has 3 differents structs, that orchestrate the execution and the concurrency of the jobs.

### Queue

Queue is the principal struct. It's the element that enqueue the execution of jobs. It has
a maxQueue max elements to control the max number of wait jobs.

* Job is a interface that provides the abstration layer to execute different tasks.
* maxQueue is a int param that set the max number of elements to queue.  
* Queue is a buffered channel that provides the input chan for jobs that will be dispathed. 

### Dispatcher

Dispatcher struct read jobs from InputQueue, and dispatch them to chanPool to be executed by
a worker. chanPool is prefetched with MaxWorkers worker that control the max concurrency 
executing jobs.

* MaxWorkers is a int param that set the max number of concurrent workers to dispatch jobs
* ChanPool is a buffered channel of channels that provides a worker chan pool. 
* InputQueue is param provided from Queue.
* Workers is a slice of worker.

### Worker

Worker struct read job from WorkerPool, and execute it. 

* WorkerPool is param provided from Dispatcher.
* WorkerChannel is a chan that execute jobs.

## How to use

'''
go get github.com/rawmind0/go-chanqueue
'''

'''
import (
    "fmt"
    "net/http"
    "github.com/rawmind0/go-chanqueue"
)
'''

At example directory there is a simple proxy service that proxy web requests to a remote url.

By default it creates a queue of 512 to accept jobs and a pool of 512 concurent workers to
execute jobs.





