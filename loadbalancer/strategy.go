package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"sync"
)

type strategy interface {
	add(handlers ...*handler)
	afterResponse(h *handler, r *http.Request)
	nextHandler(r *http.Request) *handler
}

type roundRobin struct {
	last     int
	handlers []*handler
	mu       sync.Mutex
}

func (s roundRobin) add(handlers ...*handler)                  {}
func (s roundRobin) afterResponse(h *handler, r *http.Request) {}
func (s *roundRobin) nextHandler(r *http.Request) *handler {
	s.mu.Lock()
	s.last = (s.last + 1) % len(s.handlers)
	s.mu.Unlock()
	return s.handlers[s.last%len(s.handlers)]
}

type random struct {
	handlers []*handler
}

func (s random) add(handlers []*handler)                   {}
func (s random) afterResponse(h *handler, r *http.Request) {}
func (s random) nextHandler(_ *http.Request) *handler {
	return s.handlers[rand.Intn(len(s.handlers))]
}

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

func (s leastConn) add(handlers ...*handler) {}

func (s *leastConn) afterResponse(h *handler, r *http.Request) {
	s.mu.Lock()
	s.handlers[h]--
	s.mu.Unlock()
}

func (s *leastConn) nextHandler(_ *http.Request) *handler {
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

func (s weightedRoundRobin) add(handlers ...*handler)                  {}
func (s weightedRoundRobin) afterResponse(h *handler, r *http.Request) {}
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

type simpleHashing struct {
	handlers []*handler
}

func (s simpleHashing) add(handlers ...*handler)                  {}
func (s simpleHashing) afterResponse(h *handler, r *http.Request) {}
func (s simpleHashing) nextHandler(r *http.Request) *handler {
	return s.handlers[hash(r.Header.Get("UserID"))%len(s.handlers)]
}

type consistentHashing struct {
	handlers    map[int]*handler
	logs        []string
	numReplicas int
	partitions  []int
}

type consistentHashingOpts struct {
	handlers    []*handler
	numReplicas int
}

func newConsistentHashing(opts *consistentHashingOpts) *consistentHashing {
	h := &consistentHashing{
		handlers:    map[int]*handler{},
		numReplicas: opts.numReplicas,
		partitions:  []int{},
	}

	h.add(opts.handlers...)

	return h
}

func (s *consistentHashing) add(handlers ...*handler) {
	for _, handler := range handlers {
		for i := 0; i < s.numReplicas; i++ {
			h := hash(fmt.Sprintf("%s:%d", handler.name, i))
			s.partitions = append(s.partitions, h)
			s.handlers[h] = handler
		}
	}

	sort.Ints(s.partitions)
}

func (s consistentHashing) afterResponse(h *handler, r *http.Request) {
	s.logs = append(s.logs, r.Header.Get("UserID"))
}

func (s consistentHashing) nextHandler(r *http.Request) *handler {
	hash := hash(r.Header.Get("UserID"))

	idx := sort.Search(len(s.handlers), func(i int) bool {
		return hash < s.partitions[i]
	})

	if idx == len(s.handlers) {
		idx = 0
	}

	return s.handlers[s.partitions[idx]]
}

func hash(userID string) int {
	// https://stackoverflow.com/questions/13582519/how-to-generate-hash-number-of-a-string-in-go
	// https://pkg.go.dev/hash/fnv@go1.17.3
	h := fnv.New32a()
	h.Write([]byte(userID))
	return int(h.Sum32())
}
