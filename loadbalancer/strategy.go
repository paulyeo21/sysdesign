package main

import (
	"math/rand"
	"net/http"
)

type strategy interface {
	nextHandler(r *http.Request) *handler
}

type roundRobin struct {
	last     int
	handlers []*handler
}

func (s *roundRobin) nextHandler(r *http.Request) *handler {
	s.last++
	return s.handlers[s.last%len(s.handlers)]
}

// type weightedRoundRobin struct{}

// func (s weightedRoundRobin) nextHandler(r *http.Request) *handler {
// 	// not implemented
// }

type random struct {
	handlers []*handler
}

func (s random) nextHandler(r *http.Request) *handler {
	return s.handlers[rand.Intn(len(s.handlers))]
}

// type leastLoaded struct{}

// type simpleHashing struct{}

// func (s simpleHashing) nextHandler(r *http.Request) *handler {
// 	userID := []byte(r.Header.Get("UserID"))
// 	i := int64(md5.Sum(userID))
// 	fmt.Printf("%d\n", i)
// 	return lb.handlers[i%len(lb.handlers)]
// }
