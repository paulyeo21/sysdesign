package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Handler struct {
	port string
}

func newHandler(port string) *Handler {
	return &Handler{
		port: port,
	}
}

func (h Handler) handleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from port %s", h.port)
}

func main() {
	h := newHandler(os.Args[1])
	http.HandleFunc("/", h.handleFunc)

	fmt.Printf("Listening on port %s\n", h.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", h.port), nil))
}
