package model

// Playlists is a slice of playlists received from spotify-api
type Playlists struct {
    UserID string `json:"user_id"`
    Items []Playlist `json:"items"`
}

// Playlist is a representation of a row in the playlists entity
type Playlist struct {
	ID string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	SpotifyURL string `json:"spotify_url" db:"spotify_url"`
	SpotifyID string `json:"spotify_id" db:"spotify_id"`
	TotalTracks int `json:"total_tracks" db:"total_tracks"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UserID string `json:"user_id" db:"user_id"`
}
