package main

import (
	"fmt"
	"net/url"
)

type handler struct {
	conns  int
	name   string
	url    *url.URL
	weight int
}

type handlerOpts struct {
	Name   string `json:"name"`
	Ref    string `json:"ref"`
	Weight int    `json:"weight"`
}

func newHandler(opts *handlerOpts) (*handler, error) {
	u, err := url.Parse(opts.Ref)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	return &handler{
		name:   opts.Name,
		url:    u,
		weight: opts.Weight,
	}, nil
}
