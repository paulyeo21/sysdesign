package main

type handler struct {
	cache *cache
	db    *database
}

type handlerConfig struct {
	cache *cache
	db    *database
}

func newHandler(config *handlerConfig) *handler {
	return &handler{
		cache: config.cache,
		db:    config.db,
	}
}
