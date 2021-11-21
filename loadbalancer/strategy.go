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
}

func (s *roundRobin) nextHandler(r *http.Request) *handler {
	s.last++
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

// type weightedRoundRobin struct{}

// func (s weightedRoundRobin) nextHandler(r *http.Request) *handler {
// 	// not implemented
// }

// type simpleHashing struct{}

// func (s simpleHashing) nextHandler(r *http.Request) *handler {
// 	userID := []byte(r.Header.Get("UserID"))
// 	i := int64(md5.Sum(userID))
// 	fmt.Printf("%d\n", i)
// 	return lb.handlers[i%len(lb.handlers)]
// }
