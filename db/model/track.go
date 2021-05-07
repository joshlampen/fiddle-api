package model

import (
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
)

// Tracks is a slice of tracks received from spotify-api
type Tracks struct {
    PlaylistID string  `json:"playlist_id"`
    Items      []Track `json:"items"`
}

// Track is a representation of a row in the tracks entity
type Track struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Popularity  int            `json:"popularity,omitempty" db:"popularity"`
    Duration    int            `json:"duration_ms" db:"duration"`
    AddedAt     string         `json:"added_at" db:"added_at"`
    SpotifyURI  string         `json:"spotify_uri" db:"spotify_uri"`
	SpotifyURL  string         `json:"spotify_url,omitempty" db:"spotify_url"`
	SpotifyID   string         `json:"spotify_id,omitempty" db:"spotify_id"`
    Artists     types.JSONText `json:"artists" db:"artists_json"`
    Album       types.JSONText `json:"album" db:"album_json"`
	CreatedAt   string         `json:"created_at,omitempty" db:"created_at"`
    PlaylistIDs []string       `json:"playlist_ids"`
    OwnerIDs    []string       `json:"owner_ids"`

    PlaylistIDsByteArray []byte `json:"playlist_ids_json,omitempty" db:"playlist_ids_json"`
    OwnerIDsByteArray    []byte `json:"owner_ids_json,omitempty" db:"owner_ids_json"`
}

func MapGetTracksResponse(tracks []Track) ([]Track, error) {
    var resp []Track

    for _, track := range tracks {
        var respItem Track
        respItem.ID = track.ID
        respItem.Name = track.Name
        respItem.Popularity = track.Popularity
        respItem.Duration = track.Duration
        respItem.AddedAt = track.AddedAt
        respItem.SpotifyURI = track.SpotifyURI
        respItem.Artists = track.Artists
        respItem.Album = track.Album

        playlistIDs, err := marshalBytesToStrings(track.PlaylistIDsByteArray)
        if err != nil {
            return tracks, err
        }
        respItem.PlaylistIDs = playlistIDs

        ownerIDs, err := marshalBytesToStrings(track.OwnerIDsByteArray)
        if err != nil {
            return tracks, err
        }
        respItem.OwnerIDs = ownerIDs

        resp = append(resp, respItem)
    }

    return resp, nil
}

func marshalBytesToStrings(bytes []byte) ([]string, error) {
    var resp []string

    if err := json.Unmarshal(bytes, &resp); err != nil {
        return resp, err
    }

    return resp, nil
}
