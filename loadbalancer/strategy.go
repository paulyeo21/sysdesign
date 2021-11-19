package main

type strategy interface {
	nextHandler(lb *loadBalancer) *handler
}

type roundRobin struct{}

func (s roundRobin) nextHandler(lb *loadBalancer) *handler {
	lb.latestHandler++
	return lb.handlers[lb.latestHandler%len(lb.handlers)]
}

type weightedRoundRobin struct{}

func (s weightedRoundRobin) nextHandler(lb *loadBalancer) *handler {
	// not implemented
}
