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
        "t.id AS id",
        "t.name AS name",
        "t.duration AS duration",
        "t.spotify_uri AS spotify_uri",
        "t.spotify_url AS spotify_url",
        "t.artists_json AS artists_json",
        "t.album_json AS album_json",
    }

    var filters []string
    if p.UserID != "" {
        filters = append(filters, `playlists.user_id = :user_id`)
    }

    q := `SELECT ` + strings.Join(fields, ", ") + `,
            (
                SELECT json_agg(playlists.id)
                FROM playlists
                INNER JOIN playlists_tracks ON playlists_tracks.playlist_id = playlists.id
                WHERE playlists_tracks.track_id = t.id
                GROUP BY playlists_tracks.track_id
            ) AS playlist_ids_json,
            (
                SELECT json_agg(users.id)
                FROM users
                INNER JOIN playlists ON playlists.user_id = users.id
                INNER JOIN playlists_tracks ON playlists_tracks.playlist_id = playlists.id
                WHERE playlists_tracks.track_id = t.id
                GROUP BY playlists_tracks.track_id
            ) AS owner_ids_json,
            (
                SELECT json_agg(playlists_tracks.added_at)
                FROM playlists_tracks
                WHERE playlists_tracks.track_id = t.id
                GROUP BY playlists_tracks.track_id
            ) AS added_ats_json
            FROM tracks t
            INNER JOIN playlists_tracks ON playlists_tracks.track_id = t.id
            INNER JOIN playlists ON playlists.id = playlists_tracks.playlist_id
            INNER JOIN users ON users.id = playlists.user_id
            WHERE ` + strings.Join(filters, " AND ") + `
            GROUP BY t.id, playlists_tracks.added_at
            ORDER BY playlists_tracks.added_at DESC`

    return q
}

// Search - gets row(s) in tracks based on search params
func (ts *TrackStore) Search(params TracksSearchParams) ([]model.Track, error) {
    q := params.String()

    var tracks []model.Track
    nstmt, err := ts.DB.PrepareNamed(q)
    if err != nil {
        return tracks, err
    }
    err = nstmt.Select(&tracks, params)
    if err != nil {
        return tracks, err
    }

    return tracks, nil
}

// GetBySpotifyIDIfExists - get a row from tracks by Spotify ID if it exists
func (ts *TrackStore) GetBySpotifyIDIfExists(id string) (*model.Track, error) {
    q := `SELECT EXISTS (SELECT 1 FROM tracks WHERE spotify_id = $1)`

    var exists bool
    if err := ts.DB.QueryRowx(q, id).Scan(&exists); err != nil {
        return nil, err
    }
    if !exists {
        return nil, nil
    }

    q = `SELECT * FROM tracks WHERE spotify_id = $1`

    track := &model.Track{}
    if err := ts.DB.QueryRowx(q, id).StructScan(track); err != nil {
        return nil, err
    }

    return track, nil
}

// Create - insert a row into tracks
func (ts *TrackStore) Create(t model.Track) (*model.Track, error) {
	q := `INSERT INTO tracks (
			name,
			popularity,
            duration,
            spotify_uri,
			spotify_url,
			spotify_id,
            artists_json,
            album_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *`

	track := &model.Track{}
	if err := ts.DB.QueryRowx(
		q,
		t.Name,
		t.Popularity,
        t.Duration,
        t.SpotifyURI,
		t.SpotifyURL,
		t.SpotifyID,
        t.Artists,
        t.Album,
	).StructScan(track); err != nil {
		return nil, err
	}

	return track, nil
}
