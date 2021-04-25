package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/db/model"
)

// PlaylistStore manages the playlists database entity
type PlaylistStore struct {
	DB *sqlx.DB
}

// NewPlaylistStore returns a new db connection for the PlaylistStore
func NewPlaylistStore(db *sqlx.DB) *PlaylistStore {
	return &PlaylistStore{DB: db}
}

// Create - insert a row into playlists
func (ps *PlaylistStore) Create(p model.Playlist, userID string) (model.Playlist, error) {
	q := `INSERT INTO playlists (
			name,
			spotify_url,
			spotify_id,
			total_tracks,
			user_id
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING *`

    var playlist model.Playlist
	if err := ps.DB.QueryRowx(
		q,
		p.Name,
		p.SpotifyURL,
		p.SpotifyID,
		p.TotalTracks,
		userID,
	).StructScan(&playlist); err != nil {
		return model.Playlist{}, err
	}

	return playlist, nil
}
