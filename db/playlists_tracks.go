package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/internal/utils/logger"
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
func (pts *PlaylistTrackStore) Create(playlistID, trackID, addedAt string) (string, error) {
    logger := logger.NewLogger()

	q := `INSERT INTO playlists_tracks (
			playlist_id,
			track_id,
            added_at
		) VALUES ($1, $2, $3)
		RETURNING playlist_id`

	var id []uint8
	if err := pts.DB.QueryRow(
		q,
		playlistID,
		trackID,
        addedAt,
	).Scan(&id); err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", playlistID).
            Str("trackID", trackID).
            Msg("PlaylistTrackStore.Create - failed to insert row")
		return "", err
	}

	return string(id), nil
}

// CheckExistsBySpotifyID - check if a row in playlists_tracks exists by playlist ID and track ID
func (pts *PlaylistTrackStore) CheckExistsByPlaylistIDTrackID(playlistID, trackID string) (bool, error) {
    logger := logger.NewLogger()

	q := `SELECT EXISTS (
            SELECT 1
            FROM playlists_tracks
            WHERE playlist_id = $1
            AND track_id = $2
        )`

	var exists bool
	if err := pts.DB.QueryRowx(q, playlistID, trackID).Scan(&exists); err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", playlistID).
            Str("trackID", trackID).
            Msg("PlaylistTrackStore.CheckExistsByPlaylistIDTrackID - failed to check row")
		return false, err
	}

	return exists, nil
}

// DeleteByUserID - delete all rows from playlists_tracks for a playlist ID
func (pts *PlaylistTrackStore) DeleteByPlaylistID(id string) (error) {
    logger := logger.NewLogger()

    q := `DELETE FROM playlists_tracks WHERE playlist_id = $1`

    _, err := pts.DB.Exec(q, id)
    if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", id).
            Msg("PlaylistTrackStore.DeleteByPlaylistID - failed to delete row")
        return err
    }

    return nil
}
