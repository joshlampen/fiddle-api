package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/db/model"
	jsonWriter "github.com/JoshLampen/fiddle/api/internal/utils/json"
)

// PlaylistsCreate is an HTTP handler for inserting an array of playlists into the database
func PlaylistsCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")

	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        err := fmt.Errorf("Failed to read request: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var playlists model.Playlists
	if err := json.Unmarshal(body, &playlists); err != nil {
        err := fmt.Errorf("Failed to process request: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}

    // Delete any previous playlists belonging to the user
    err = store.PlaylistStore.DeleteByUserID(playlists.UserID)
    if err != nil {
        err := fmt.Errorf("Failed to clear playlists before creating: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
        return
    }

	// Insert into database
	for i, playlist := range playlists.Items {
        result, err := store.PlaylistStore.Create(playlist, playlists.UserID)
		if err != nil {
            err := fmt.Errorf("Failed to create playlist: %w", err)
            jsonWriter.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		playlists.Items[i].ID = result.ID
	}

    jsonWriter.WriteResponse(w, playlists)
}
