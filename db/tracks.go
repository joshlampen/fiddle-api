package db

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/db/model"
)

// TrackStore manages the tracks database entity
type TrackStore struct {
	DB *sqlx.DB
}

// NewTrackStore returns a new db connection for the TrackStore
func NewTrackStore(db *sqlx.DB) *TrackStore {
	return &TrackStore{DB: db}
}

type TracksSearchParams struct {
    UserID string `db:"user_id"`
}

func (p TracksSearchParams) String() string {
    fields := []string {
        "tracks.id as id",
        "tracks.name as name",
        "tracks.duration as duration",
        "tracks.spotify_uri as spotify_uri",
        "tracks.spotify_url as spotify_url",
        "tracks.artists_json as artists_json",
        "tracks.album_json as album_json",
        "playlists.name as playlist_name",
        "playlists.spotify_url as playlist_spotify_url",
    }

    var filters []string
    if p.UserID != "" {
        filters = append(filters, `playlists.user_id = :user_id`)
    }

    q := `SELECT ` + strings.Join(fields, ", ") + `
            FROM tracks
            INNER JOIN playlists_tracks on playlists_tracks.track_id = tracks.id
            INNER JOIN playlists on playlists.id = playlists_tracks.playlist_id
            WHERE ` + strings.Join(filters, " AND ") + `
            ORDER BY tracks.added_at DESC`

    return q
}

// Search - gets row(s) in tracks based on search params
func (ts *TrackStore) Search(params TracksSearchParams) ([]model.Track, error) {
    q := params.String()

    var tracks []model.Track
    nstmt, err := ts.DB.PrepareNamed(q)
    if err != nil {
        return []model.Track{}, err
    }
    err = nstmt.Select(&tracks, params)
    if err != nil {
        return []model.Track{}, err
    }

    return tracks, nil
}

// GetBySpotifyID - gets a row in tracks by Spotify ID
func (ts *TrackStore) GetBySpotifyID(id string) (model.Track, error) {
	q := `SELECT * FROM tracks WHERE spotify_id = $1`

	var track model.Track
	if err := ts.DB.QueryRowx(q, id).StructScan(&track); err != nil {
		return model.Track{}, err
	}

	return track, nil
}

// Create - insert a row into tracks
func (ts *TrackStore) Create(t model.Track) (model.Track, error) {
	q := `INSERT INTO tracks (
			name,
			popularity,
            duration,
            added_at,
            spotify_uri,
			spotify_url,
			spotify_id,
            artists_json,
            album_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING *`

	var track model.Track
	if err := ts.DB.QueryRowx(
		q,
		t.Name,
		t.Popularity,
        t.Duration,
        t.AddedAt,
        t.SpotifyURI,
		t.SpotifyURL,
		t.SpotifyID,
        t.Artists,
        t.Album,
	).StructScan(&track); err != nil {
		return model.Track{}, err
	}

	return track, nil
}

// CheckExistsBySpotifyID - check if a row in tracks exists by Spotify ID
func (ts *TrackStore) CheckExistsBySpotifyID(id string) (bool, error) {
	q := `SELECT EXISTS (SELECT 1 FROM tracks WHERE spotify_id = $1)`

	var exists bool
	if err := ts.DB.QueryRowx(q, id).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
