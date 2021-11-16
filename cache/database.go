package main

import "time"

type database struct {
	values map[string]string
}

func newDatabase() *database {
	values := map[string]string{
		"index": "Hi there, I love monkeys",
	}

	return &database{
		values: values,
	}
}

func (d database) get(key string) string {
	time.Sleep(3 * time.Second)

	return d.values[key]
}
