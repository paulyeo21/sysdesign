package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

var validPath = regexp.MustCompile("^/(nocache|cache)$")

// cache
func (h handler) cachedHandler(w http.ResponseWriter, r *http.Request) {
	val, err := h.cache.get("index")
	if err != nil {
		val = h.db.get("index")
		h.cache.set("index", val)
	}

	fmt.Fprintf(w, val)
}

// nocache
func (h handler) noCacheHandler(w http.ResponseWriter, r *http.Request) {
	val := h.db.get("index")
	fmt.Fprintf(w, val)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r)

		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("%v\n", elapsed)
	}
}

func main() {
	handler := newHandler(&handlerConfig{
		cache: newCache(),
		db:    newDatabase(),
	})

	http.HandleFunc("/cache", makeHandler(handler.cachedHandler))
	http.HandleFunc("/nocache", makeHandler(handler.noCacheHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
