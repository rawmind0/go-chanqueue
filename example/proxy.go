package main

import (
    "fmt"
    "net/http"
    "github.com/rawmind0/go-chanqueue"
)

type Proxy struct {
	Listen string
	Queue *chanqueue.Queue
	Url string
	Quit chan bool
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	var j chanqueue.Job
	j = NewSite(w, r, p.Url)
	go p.Queue.EnqueueJob(j)
	j.WaitJob()
}

func NewProxy(l string, maxQueue int, maxWorkers int, url string) *Proxy {
	return &Proxy{
		Listen: l, 
		Queue: chanqueue.NewQueue(maxQueue, maxWorkers),
		Url: url
		Quit: make(chan bool)}
}

func (p *Proxy) Start() {
	p.Queue.Start()
	p.StartEndpoints()
    http.HandleFunc("/", p.handler)
	fmt.Printf("Listenning at: %s\n", p.Listen)
    http.ListenAndServe(p.Listen, nil)
}