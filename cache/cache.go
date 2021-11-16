package main

import "errors"

type cache struct {
	values map[string]string
}

func newCache() *cache {
	return &cache{
		values: map[string]string{},
	}
}

func (c *cache) set(key, value string) {
	c.values[key] = value
}

func (c cache) get(key string) (string, error) {
	res, ok := c.values[key]

	if !ok {
		return "", errors.New("not found")
	}

	return res, nil
}
