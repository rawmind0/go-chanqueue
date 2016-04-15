package main

import (
	"fmt"
	"net/http"
	"net"
	"time"
	"sync"
	"io/ioutil"
)

type Site struct {
	Writer http.ResponseWriter
	Reader *http.Request
	Response *http.Response
	Quit chan bool
	Status string
	Mutex *sync.Mutex
	TimeResponse time.Duration
}

func NewSite(w http.ResponseWriter, r *http.Request, url string) *Site {
	return &Site{
		Writer: w, 
		Reader: r, 
		Response: &http.Response{},
		Url: url
		Quit: make(chan bool)}
}

func (s *Site) Start() {
	r,t,_:=DoRequest(s.Reader.Method, s.Url)	
	fmt.Printf("[%s] %s://%s%s From: %s Took:%s\n", s.Reader.Method, s.Reader.URL.Scheme, s.Url, s.Reader.URL.Path, s.Reader.RemoteAddr, t)
	fmt.Fprintf(s.Writer, "%s\n", r)
	s.Stop()
}

func (s *Site) Stop() {
	s.Quit <- true
}

func (s *Site) Wait() {
	<- s.Quit
}

func NewClient () (c *http.Client){

	var secs time.Duration
	secs = 3 // rather aggressive
	c = &http.Client{
  		Transport: &http.Transport{
    		Proxy: http.ProxyFromEnvironment,
    		Dial: (&net.Dialer{
      			Timeout:   secs * time.Second,
      			KeepAlive: 5 * time.Second,
    		}).Dial,
    		TLSHandshakeTimeout: secs * time.Second,
    		MaxIdleConnsPerHost: 512,
    		DisableKeepAlives: true,
  		},
	}

	return 
}

func DoRequest (method string,url string) (string, time.Duration, error) {
	start := time.Now()
	var s string
	client := NewClient()
	req , err := http.NewRequest(method, url, nil)

	if err == nil {
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()

			contents, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				s = string(contents)
			} else {
    			fmt.Println(err)
    			s = "ERROR"
			}
		} else {
			fmt.Println(err)
    		s = "ERROR"
		}
	}else {
		fmt.Println(err)
    	s = "ERROR"
	}
	elapsed := time.Since(start)
	return s,elapsed,err

}