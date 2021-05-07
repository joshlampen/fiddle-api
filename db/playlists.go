package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/db/model"
	"github.com/JoshLampen/fiddle/api/internal/utils/logger"
)

// PlaylistStore manages the playlists database entity
type PlaylistStore struct {
	DB *sqlx.DB
}

// NewPlaylistStore returns a new db connection for the PlaylistStore
func NewPlaylistStore(db *sqlx.DB) *PlaylistStore {
	return &PlaylistStore{DB: db}
}

// GetByUserID - get a user and their friends' playlists by user ID
func (ps *PlaylistStore) GetByUserID(id string) ([]model.Playlist, error) {
    logger := logger.NewLogger()

    q := `SELECT * FROM playlists
            WHERE playlists.user_id = :user_id
            UNION
            SELECT p.* FROM playlists p
            INNER JOIN friendships ON friendships.friend_id = p.user_id
            WHERE friendships.user_id = :user_id
            AND friendships.pending = false`

    params := struct {
        UserID string `db:"user_id"`
    } {
        UserID: id,
    }

    var playlists []model.Playlist
    nstmt, err := ps.DB.PrepareNamed(q)
    if err != nil {
        logger.Error().
            Err(err).
            Str("userID", id).
            Msg("PlaylistStore.GetByUserID - failed to prepare named statement")
        return playlists, err
    }
    err = nstmt.Select(&playlists, params)
    if err != nil {
        logger.Error().
            Err(err).
            Str("userID", id).
            Msg("PlaylistStore.GetByUserID - failed to get playlists")
        return playlists, err
    }

    return playlists, nil
}

// Create - insert a row into playlists
func (ps *PlaylistStore) Create(p model.Playlist, userID string) (model.Playlist, error) {
    logger := logger.NewLogger()

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
        logger.Error().
            Err(err).
            Str("playlistID", p.ID).
            Str("userID", userID).
            Msg("PlaylistStore.Create - failed to create playlist")
		return model.Playlist{}, err
	}

	return playlist, nil
}

// DeleteByUserID - delete all rows from playlists for a user ID
func (ps *PlaylistStore) DeleteByUserID(id string) (error) {
    logger := logger.NewLogger()

    q := `DELETE FROM playlists WHERE user_id = $1`

    _, err := ps.DB.Exec(q, id)
    if err != nil {
        logger.Error().
            Err(err).
            Str("userID", id).
            Msg("PlaylistStore.DeleteByUserID - failed to delete playlist")
        return err
    }

    return nil
}
