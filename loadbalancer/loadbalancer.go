package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

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

	// strategy callback after res
	lb.strategy.afterResponse(h)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fn(w, r)

		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("%v\n", elapsed)
	}
}

func main() {
	// lb := newLoadBalancer(
	// 	&roundRobin{
	// 		handlers: []*handler{
	// 			newHandler("http://localhost:3000"),
	// 			newHandler("http://localhost:3001"),
	// 			newHandler("http://localhost:3002"),
	// 		},
	// 	},
	// )

	// lb := newLoadBalancer(
	// 	newWeightedRoundRobin([]*handler{
	// 		newHandler(&handlerOpts{
	// 			ref:    "http://localhost:3000",
	// 			weight: 3,
	// 		}),
	// 		newHandler(&handlerOpts{
	// 			ref:    "http://localhost:3001",
	// 			weight: 1,
	// 		}),
	// 		newHandler(&handlerOpts{
	// 			ref:    "http://localhost:3002",
	// 			weight: 1,
	// 		}),
	// 	}),
	// )

	// lb := newLoadBalancer(
	// 	newLeastConn([]*handler{
	// 		newHandler("http://localhost:3000"),
	// 		newHandler("http://localhost:3001"),
	// 		newHandler("http://localhost:3002"),
	// 	}),
	// )

	lb := newLoadBalancer(simpleHashing{
		hs: []*handler{
			newHandler(&handlerOpts{ref: "http://localhost:3000"}),
			newHandler(&handlerOpts{ref: "http://localhost:3001"}),
			newHandler(&handlerOpts{ref: "http://localhost:3002"}),
		},
	})

	http.HandleFunc("/", makeHandler(lb.handleFunc))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
