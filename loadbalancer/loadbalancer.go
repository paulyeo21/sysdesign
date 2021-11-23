package main

import (
	"encoding/json"
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
	lb.strategy.afterResponse(h, r)
}

func (lb *loadBalancer) addHandleFunc(w http.ResponseWriter, r *http.Request) {
	// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
	dec := json.NewDecoder(r.Body)

	var h handlerOpts

	err := dec.Decode(&h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handler, err := newHandler(&h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lb.strategy.add(handler)
	fmt.Fprintf(w, "Added handler at %v\n", handler.url.Host)
}

func (lb *loadBalancer) reportHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", lb.strategy.report())
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fn(w, r)

		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Println(elapsed)
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

	// lb := newLoadBalancer(simpleHashing{
	// 	hs: []*handler{
	// 		newHandler(&handlerOpts{ref: "http://localhost:3000"}),
	// 		newHandler(&handlerOpts{ref: "http://localhost:3001"}),
	// 		newHandler(&handlerOpts{ref: "http://localhost:3002"}),
	// 	},
	// })

	h1, err := newHandler(&handlerOpts{
		Name: "node1",
		Ref:  "http://localhost:3000",
	})
	if err != nil {
		log.Fatal(err)
	}

	h2, err := newHandler(&handlerOpts{
		Name: "node2",
		Ref:  "http://localhost:3001",
	})
	if err != nil {
		log.Fatal(err)
	}

	h3, err := newHandler(&handlerOpts{
		Name: "node3",
		Ref:  "http://localhost:3002",
	})
	if err != nil {
		log.Fatal(err)
	}

	lb := newLoadBalancer(newConsistentHashing(
		&consistentHashingOpts{
			handlers:    []*handler{h1, h2, h3},
			numReplicas: 10,
		},
	))

	http.HandleFunc("/", makeHandler(lb.handleFunc))
	http.HandleFunc("/add", lb.addHandleFunc)
	http.HandleFunc("/report", lb.reportHandleFunc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
