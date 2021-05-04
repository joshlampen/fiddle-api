package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001" // localhost
	}

	dbConn, err := db.Init(port)
	if err != nil {
		err := fmt.Errorf("failed to initialize DB: %w", err)
		panic(err)
	}
	defer dbConn.Close()
	store := db.NewStore(dbConn)

	r := mux.NewRouter()

	r.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) { handler.TokensGet(w, r, store) }).Methods("GET", "OPTIONS")
	r.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) { handler.TokensCreate(w, r, store) }).Methods("POST", "OPTIONS")

	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) { handler.UsersGet(w, r, store) }).Methods("GET", "OPTIONS")
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) { handler.UsersCreate(w, r, store) }).Methods("POST", "OPTIONS")

	r.HandleFunc("/playlists", func(w http.ResponseWriter, r *http.Request) { handler.PlaylistsCreate(w, r, store) }).Methods("POST", "OPTIONS")

	r.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) { handler.TracksSearch(w, r, store) }).Methods("GET", "OPTIONS")
	r.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) { handler.TracksCreate(w, r, store) }).Methods("POST", "OPTIONS")

	http.ListenAndServe(":" + port, r)
}
