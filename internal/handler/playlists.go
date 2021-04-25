package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/db/model"
)

// PlaylistsCreate is an HTTP handler for inserting an array of playlists into the database
func PlaylistsCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")

	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("handler.PlaylistsCreate - failed to read request body:", err)
		return
	}

	var playlists model.Playlists
	if err := json.Unmarshal(body, &playlists); err != nil {
		fmt.Println("handler.PlaylistsCreate - failed to unmarshal request body:", err)
		return
	}

	// Insert into database
	for i, playlist := range playlists.Items {
		result, err := store.PlaylistStore.Create(playlist, playlists.UserID)
		if err != nil {
			fmt.Println("handler.PlaylistsCreate - failed to create playlist:", err)
			return
		}
		playlists.Items[i].ID = result.ID
	}

	// Send a response
	jsonBody, err := json.Marshal(playlists)
	if err != nil {
		fmt.Println("handler.PlaylistsCreate - failed to marshal response body:", err)
		return
	}
	w.Write(jsonBody)
}
