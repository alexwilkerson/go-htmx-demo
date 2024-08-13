package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexwilkerson/go-htmx-demo/internal/components"
	petname "github.com/dustinkirkland/golang-petname"
)

func main() {
	count := 0
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	leaderboard := map[string]int{}
	chat := []components.Message{}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		user := sessionManager.GetString(r.Context(), "user")
		if user == "" {
			user = petname.Generate(2, "-")
			sessionManager.Put(r.Context(), "user", user)
		}

		components.Index(user).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /chat", func(w http.ResponseWriter, r *http.Request) {
		const chatSize = 10

		window := r.URL.Query().Get("window")

		if window == strings.ToLower("true") {
			if len(chat) > chatSize {
				chat = chat[len(chat)-chatSize:]
			}

			components.ChatWindow(chat).Render(r.Context(), w)
			return
		}

		components.Chat().Render(r.Context(), w)
	})

	mux.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {
		user := sessionManager.GetString(r.Context(), "user")
		message := r.FormValue("message")

		if message == "" {
			return
		}

		chat = append(chat, components.Message{
			User:    user,
			Message: message,
			Time:    time.Now(),
		})
	})

	mux.HandleFunc("GET /leaderboard", func(w http.ResponseWriter, r *http.Request) {
		components.Leaderboard(leaderboard).Render(r.Context(), w)
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

		leaderboard[user] = sessionCount + 1

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
