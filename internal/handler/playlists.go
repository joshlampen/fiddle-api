package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/db/model"
	"github.com/JoshLampen/fiddle/api/internal/constant"
	jsonWriter "github.com/JoshLampen/fiddle/api/internal/utils/json"
)

func PlaylistsGetByUserID(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")

    userID := r.URL.Query().Get(constant.URLParamUserID)

    // Retrieve from database
    playlists, err := store.PlaylistStore.GetByUserID(userID)
    if err != nil {
        err := fmt.Errorf("Failed to get playlists: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
        return
    }

    jsonWriter.WriteResponse(w, playlists)
}

// PlaylistsCreate is an HTTP handler for inserting an array of playlists into the database
func PlaylistsCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")

	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        err := fmt.Errorf("Failed to read playlists request: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var playlists model.Playlists
	if err := json.Unmarshal(body, &playlists); err != nil {
        err := fmt.Errorf("Failed to process playlists request: %w", err)
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
