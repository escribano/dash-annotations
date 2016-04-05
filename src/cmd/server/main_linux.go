// +build !darwin
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Listening on %q", listen)
	http.ListenAndServe(listen, logHandler(jsonHandler(mux)))
}
