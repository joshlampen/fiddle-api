package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/db/model"
)

func TracksSearch(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")

    // Read the request
    var params db.TracksSearchParams

    if r.FormValue("user_id") != "" {
        params.UserID = r.FormValue("user_id")
    }

    // Retrieve from database
    tracks, err := store.TrackStore.Search(params)
    if err != nil {
        fmt.Println("handler.TracksSearch - failed to get tracks:", err)
        return
    }

    // Send a response
	jsonBody, err := json.Marshal(tracks)
	if err != nil {
		fmt.Println("handler.TracksSearch - failed to marshal response body:", err)
		return
	}
	w.Write(jsonBody)
}

// TracksCreate is an HTTP handler for inserting an array of tracks into the database
func TracksCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")
	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("handler.TracksCreate - failed to read request body:", err)
		return
	}

	var tracks model.Tracks
	if err := json.Unmarshal(body, &tracks); err != nil {
		fmt.Println("handler.TracksCreate - failed to unmarshal request body:", err)
		return
	}

	// Insert into database
    var respBody model.Tracks
	for _, item := range tracks.Items {
        // Check if the track exists
		tExists, err := store.TrackStore.CheckExistsBySpotifyID(item.SpotifyID)
		if err != nil {
			fmt.Println("handler.TracksCreate - failed to check if track exists:", err)
			return
		}

		var track model.Track
		if tExists {
			result, err := store.TrackStore.GetBySpotifyID(item.SpotifyID)
			if err != nil {
				fmt.Println("handler.TracksCreate - failed to get track:", err)
				return
			}
			track = result
		} else {
			result, err := store.TrackStore.Create(item)
			if err != nil {
				fmt.Println("handler.TracksCreate - failed to create track:", err)
				return
			}
			track = result
            respBody.Items = append(respBody.Items, track)
		}

        // Check if the playlists_tracks relationship exists
        // (avoiding track duplicates in playlists)
        ptExists, err := store.PlaylistTrackStore.CheckExistsByPlaylistIDTrackID(item.PlaylistID, track.ID)
        if err != nil {
			fmt.Println("handler.TracksCreate - failed to check if playlists_tracks relationship exists:", err)
			return
		}

        if !ptExists {
            _, err = store.PlaylistTrackStore.Create(item.PlaylistID, track.ID)
            if err != nil {
                fmt.Println("handler.TracksCreate - failed to create playlists_tracks relationship:", err)
                return
            }
        }
	}

    // Send a response
	jsonBody, err := json.Marshal(respBody)
	if err != nil {
		fmt.Println("handler.TracksCreate - failed to marshal response body:", err)
		return
	}
	w.Write(jsonBody)
}
