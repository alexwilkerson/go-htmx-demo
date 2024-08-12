package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	count := 0

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.html")
	})

	http.HandleFunc("GET /button", func(w http.ResponseWriter, r *http.Request) {
		count++

		w.Write([]byte(fmt.Sprintf("<p>Count: %d</p>", count)))
	})

	http.Handle("GET /static/*",
		http.StripPrefix("/static/",
			http.FileServer(
				http.Dir("./static"),
			),
		),
	)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
