package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type handler struct {
	url *url.URL
}

func newHandler(ref string) *handler {
	u, _ := url.Parse(ref)

	return &handler{
		url: u,
	}
}

type loadBalancer struct {
	handlers      []*handler
	latestHandler int
}

func newLoadBalancer() *loadBalancer {
	hA := newHandler("http://localhost:3000")
	hB := newHandler("http://localhost:3001")
	hC := newHandler("http://localhost:3002")

	return &loadBalancer{
		handlers: []*handler{hA, hB, hC},
	}
}

func (lb *loadBalancer) handleFunc(w http.ResponseWriter, r *http.Request) {
	h := lb.nextHandler()
	rp := httputil.NewSingleHostReverseProxy(h.url)
	rp.ServeHTTP(w, r)
}

// round robin
func (lb *loadBalancer) nextHandler() *handler {
	lb.latestHandler++
	return lb.handlers[lb.latestHandler%len(lb.handlers)]
}

func main() {
	lb := newLoadBalancer()
	http.HandleFunc("/", lb.handleFunc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
