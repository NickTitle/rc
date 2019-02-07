package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.Int("p", 3000, "port to bind to")
	flag.Parse()

	http.Handle("/", blackhole())
	http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
}

func blackhole() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
