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
	strategy strategy
}

func newLoadBalancer(s strategy) *loadBalancer {
	return &loadBalancer{
		strategy: s,
	}
}

func (lb *loadBalancer) handleFunc(w http.ResponseWriter, r *http.Request) {
	h := lb.strategy.nextHandler(r)
	rp := httputil.NewSingleHostReverseProxy(h.url)
	rp.ServeHTTP(w, r)
}

func main() {
	lb := newLoadBalancer(random{
		handlers: []*handler{
			newHandler("http://localhost:3000"),
			newHandler("http://localhost:3001"),
			newHandler("http://localhost:3002"),
		},
	})

	http.HandleFunc("/", lb.handleFunc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
