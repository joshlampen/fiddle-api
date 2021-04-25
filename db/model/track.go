package model

import "github.com/jmoiron/sqlx/types"

// Tracks is a slice of tracks received from spotify-api
type Tracks struct {
    Items []Track `json:"items"`
}

// Track is a representation of a row in the tracks entity
type Track struct {
	ID string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Popularity int `json:"popularity,omitempty" db:"popularity"`
    Duration int `json:"duration_ms" db:"duration"`
    AddedAt string `json:"added_at" db:"added_at"`
    SpotifyURI string `json:"spotify_uri" db:"spotify_uri"`
	SpotifyURL string `json:"spotify_url" db:"spotify_url"`
	SpotifyID string `json:"spotify_id,omitempty" db:"spotify_id"`
    Artists types.JSONText `json:"artists" db:"artists_json"`
    Album types.JSONText `json:"album" db:"album_json"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
    PlaylistName string `json:"playlist_name" db:"playlist_name"`
    PlaylistSpotifyURL string `json:"playlist_spotify_url" db:"playlist_spotify_url"`
	PlaylistID string `json:"playlist_id,omitempty"`
}
