package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, New York Times!"))
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
