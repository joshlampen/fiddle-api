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

func TracksSearch(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")

    // Read the request
    var params db.TracksSearchParams

    if r.FormValue(constant.URLParamUserID) != "" {
        params.UserID = r.FormValue(constant.URLParamUserID)
    }

    // Retrieve from database
    tracks, err := store.TrackStore.Search(params)
    if err != nil {
        err := fmt.Errorf("Failed to get tracks: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
        return
    }

    resp, err := model.MapGetTracksResponse(tracks)
    if err != nil {
        err := fmt.Errorf("Failed to prepare tracks response: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
        return
    }

    jsonWriter.WriteResponse(w, resp)
}

// TracksCreate is an HTTP handler for inserting an array of tracks into the database
func TracksCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")
	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        err := fmt.Errorf("Failed to read tracks request: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var tracks model.Tracks
	if err := json.Unmarshal(body, &tracks); err != nil {
        err := fmt.Errorf("Failed to process tracks request: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}

    // Delete any previous playlists_tracks relationships belonging to the playlist
    err = store.PlaylistTrackStore.DeleteByPlaylistID(tracks.PlaylistID)
    if err != nil {
        err := fmt.Errorf("Failed to clear playlist tracks before creating: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
        return
    }

	// Insert into database
    var respBody model.Tracks
	for _, item := range tracks.Items {
        // Use track from the database if it exists
        track, err := store.TrackStore.GetBySpotifyIDIfExists(item.SpotifyID)
        if err != nil {
            err := fmt.Errorf("Failed to get track: %w", err)
            jsonWriter.WriteError(w, err, http.StatusInternalServerError)
            return
        }
        if track == nil {
            // Insert into database
			track, err = store.TrackStore.Create(item)
			if err != nil {
                err := fmt.Errorf("Failed to create track: %w", err)
                jsonWriter.WriteError(w, err, http.StatusInternalServerError)
				return
			}
        }
        respBody.Items = append(respBody.Items, *track)

        // Check if the playlists_tracks relationship exists
        // (avoiding track duplicates in playlists)
        ptExists, err := store.PlaylistTrackStore.CheckExistsByPlaylistIDTrackID(tracks.PlaylistID, track.ID)
        if err != nil {
            err := fmt.Errorf("Failed to check playlist: %w", err)
            jsonWriter.WriteError(w, err, http.StatusInternalServerError)
			return
		}

        if !ptExists {
            _, err = store.PlaylistTrackStore.Create(tracks.PlaylistID, track.ID, item.AddedAt)
            if err != nil {
                err := fmt.Errorf("Failed to create tracks in playlist: %w", err)
                jsonWriter.WriteError(w, err, http.StatusInternalServerError)
                return
            }
        }
	}

    jsonWriter.WriteResponse(w, respBody)
}
