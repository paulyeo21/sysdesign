package main

import (
	"math"
	"math/rand"
	"net/http"
	"sync"
)

type strategy interface {
	nextHandler(r *http.Request) *handler
	afterResponse(h *handler)
}

type roundRobin struct {
	last     int
	handlers []*handler
	mu       sync.Mutex
}

func (s *roundRobin) nextHandler(r *http.Request) *handler {
	s.mu.Lock()
	s.last = (s.last + 1) % len(s.handlers)
	s.mu.Unlock()
	return s.handlers[s.last%len(s.handlers)]
}

func (s roundRobin) afterResponse(h *handler) {}

type random struct {
	handlers []*handler
}

func (s random) nextHandler(r *http.Request) *handler {
	return s.handlers[rand.Intn(len(s.handlers))]
}

func (s random) afterResponse(h *handler) {}

type leastConn struct {
	mu       sync.Mutex
	handlers map[*handler]int
}

func newLeastConn(handlers []*handler) *leastConn {
	hs := map[*handler]int{}

	for i := range handlers {
		hs[handlers[i]] = 0
	}

	return &leastConn{
		handlers: hs,
	}
}

func (s *leastConn) nextHandler(r *http.Request) *handler {
	var res *handler
	min := math.MaxInt32

	s.mu.Lock()
	for k, v := range s.handlers {
		if v < min {
			res = k
			min = v
		}
	}

	s.handlers[res]++
	s.mu.Unlock()

	return res
}

func (s *leastConn) afterResponse(h *handler) {
	s.mu.Lock()
	s.handlers[h]--
	s.mu.Unlock()
}

type weightedRoundRobin struct {
	hs []*handler
	ws []int
}

func newWeightedRoundRobin(handlers []*handler) *weightedRoundRobin {
	var sum int
	weights := []int{}

	for i := range handlers {
		sum += handlers[i].weight
		weights = append(weights, sum)
	}

	return &weightedRoundRobin{
		hs: handlers,
		ws: weights,
	}
}

func (s weightedRoundRobin) nextHandler(_ *http.Request) *handler {
	randInt := rand.Intn(s.ws[len(s.ws)-1] + 1) // rand.Intn is [0, n)

	var mid int
	l, r := 0, len(s.ws)-1

	for l < r {
		mid = l + (r-l)/2

		if randInt <= s.ws[mid] {
			r = mid
		} else {
			l = mid + 1
		}
	}

	return s.hs[l]
}

func (s weightedRoundRobin) afterResponse(h *handler) {}

// type simpleHashing struct{}

// func (s simpleHashing) nextHandler(r *http.Request) *handler {
// 	userID := []byte(r.Header.Get("UserID"))
// 	i := int64(md5.Sum(userID))
// 	fmt.Printf("%d\n", i)
// 	return lb.handlers[i%len(lb.handlers)]
// }
