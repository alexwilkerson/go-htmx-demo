package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexwilkerson/go-htmx-demo/internal/components"
	petname "github.com/dustinkirkland/golang-petname"
)

func main() {
	count := 0
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		user := sessionManager.GetString(r.Context(), "user")
		if user == "" {
			user = petname.Generate(2, "-")
			sessionManager.Put(r.Context(), "user", user)
		}

		components.Index(user).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /button", func(w http.ResponseWriter, r *http.Request) {
		user := sessionManager.GetString(r.Context(), "user")
		if user == "" {
			user = petname.Generate(2, "-")
			sessionManager.Put(r.Context(), "user", user)
		}

		count++

		sessionCount := sessionManager.GetInt(r.Context(), "count")
		sessionManager.Put(r.Context(), "count", sessionCount+1)

		w.Write([]byte(fmt.Sprintf("<p>User: %s</p><p>Global Count: %d</p><p>Session Count: %d</p>", user, count, sessionCount+1)))
	})

	mux.Handle("GET /static/*",
		http.StripPrefix("/static/",
			http.FileServer(
				http.Dir("./static"),
			),
		),
	)

	log.Fatal(http.ListenAndServe(":3000", sessionManager.LoadAndSave(mux)))
}
