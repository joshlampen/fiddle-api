package model

// User is a representation of a row in the users entity
type User struct {
	ID string `json:"id" db:"id"`
	DisplayName string `json:"display_name" db:"display_name"`
	Email string `json:"email" db:"email"`
	SpotifyURL string `json:"spotify_url" db:"spotify_url"`
	SpotifyImageURL string `json:"spotify_image_url" db:"spotify_image_url"`
	SpotifyID string `json:"spotify_id" db:"spotify_id"`
    Token string `json:"token" db:"token"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
