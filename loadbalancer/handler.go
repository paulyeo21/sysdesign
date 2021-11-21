package main

import "net/url"

type handler struct {
	conns int
	url   *url.URL
}

func newHandler(ref string) *handler {
	u, _ := url.Parse(ref)

	return &handler{
		url: u,
	}
}
