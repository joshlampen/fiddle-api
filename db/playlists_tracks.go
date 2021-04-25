package db

import (
	"github.com/jmoiron/sqlx"
)

// PlaylistTrackStore manages the playlists_tracks database entity
type PlaylistTrackStore struct {
	DB *sqlx.DB
}

// NewPlaylistTrackStore returns a new db connection for the PlaylistTrackStore
func NewPlaylistTrackStore(db *sqlx.DB) *PlaylistTrackStore {
	return &PlaylistTrackStore{DB: db}
}

// Create - insert a row into playlists_tracks
func (pts *PlaylistTrackStore) Create(playlistID, trackID string) (string, error) {
	q := `INSERT INTO playlists_tracks (
			playlist_id,
			track_id
		) VALUES ($1, $2)
		RETURNING playlist_id`

	var id []uint8
	if err := pts.DB.QueryRow(
		q,
		playlistID,
		trackID,
	).Scan(&id); err != nil {
		return "", err
	}

	return string(id), nil
}

// CheckExistsBySpotifyID - check if a row in playlists_tracks exists by playlist ID and track ID
func (pts *PlaylistTrackStore) CheckExistsByPlaylistIDTrackID(playlistID, trackID string) (bool, error) {
	q := `SELECT EXISTS (
            SELECT 1
            FROM playlists_tracks
            WHERE playlist_id = $1
            AND track_id = $2
        )`

	var exists bool
	if err := pts.DB.QueryRowx(q, playlistID, trackID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
