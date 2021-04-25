package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/internal/handler"
)

func main() {
	dbConn, err := db.Init()
	if err != nil {
		fmt.Println("Failed to initialize DB:", err)
		panic(err)
	}
	defer dbConn.Close()
	store := db.NewStore(dbConn)

	r := mux.NewRouter()

	r.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) { handler.TokensGet(w, r, store) }).Methods("GET", "OPTIONS")
	r.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) { handler.TokensCreate(w, r, store) }).Methods("POST", "OPTIONS")

	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) { handler.UsersCreate(w, r, store) }).Methods("POST", "OPTIONS")
	r.HandleFunc("/playlists", func(w http.ResponseWriter, r *http.Request) { handler.PlaylistsCreate(w, r, store) }).Methods("POST", "OPTIONS")

	r.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) { handler.TracksSearch(w, r, store) }).Methods("GET", "OPTIONS")
	r.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) { handler.TracksCreate(w, r, store) }).Methods("POST", "OPTIONS")

	http.ListenAndServe(":8080", r)
}