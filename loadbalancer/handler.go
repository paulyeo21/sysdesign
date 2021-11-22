package main

import "net/url"

type handler struct {
	conns  int
	url    *url.URL
	weight int
}

type handlerOpts struct {
	ref    string
	weight int
}

func newHandler(opts *handlerOpts) *handler {
	u, _ := url.Parse(opts.ref)

	return &handler{
		url:    u,
		weight: opts.weight,
	}
}
