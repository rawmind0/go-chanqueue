package main

import (
	"flag"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		workers = flag.Int("workers", 512, "Number of concurrent workers")
		queues = flag.Int("queues", 512, "Size of jobs queue")
		url = flag.Int("url", "http://localhost:8090", "url where to proxy request")
	)
	flag.Parse()

	NewProxy(*listen, *queues, *workers, *url.Start()

}