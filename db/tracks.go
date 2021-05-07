package db

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/db/model"
	"github.com/JoshLampen/fiddle/api/internal/utils/logger"
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

    // var filters []string
    // if p.UserID != "" {
    //     filters = append(filters, `playlists.user_id = :user_id`)
    // }

    q := `SELECT ` + strings.Join(fields, ", ") + `,
            (
                SELECT playlists_tracks.added_at
                FROM playlists_tracks
                WHERE playlists_tracks.track_id = t.id
                ORDER BY playlists_tracks.added_at DESC
                LIMIT 1
            ) AS added_at,
            (
                SELECT jsonb_agg(playlists.id ORDER BY playlists_tracks.added_at ASC)
                FROM playlists
                INNER JOIN playlists_tracks ON playlists_tracks.playlist_id = playlists.id
                WHERE playlists_tracks.track_id = t.id
                GROUP BY playlists_tracks.track_id
            ) AS playlist_ids_json,
            (
                SELECT jsonb_agg(users.id ORDER BY playlists_tracks.added_at ASC)
                FROM users
                INNER JOIN playlists ON playlists.user_id = users.id
                INNER JOIN playlists_tracks ON playlists_tracks.playlist_id = playlists.id
                WHERE playlists_tracks.track_id = t.id
                GROUP BY playlists_tracks.track_id
            ) AS owner_ids_json
            FROM (
                SELECT t1.*
                FROM tracks t1
                INNER JOIN playlists_tracks ON playlists_tracks.track_id = t1.id
                INNER JOIN playlists ON playlists.id = playlists_tracks.playlist_id
                INNER JOIN users ON users.id = playlists.user_id
                WHERE users.id = :user_id
                UNION
                SELECT t2.*
                FROM tracks t2
                INNER JOIN playlists_tracks ON playlists_tracks.track_id = t2.id
                INNER JOIN playlists ON playlists.id = playlists_tracks.playlist_id
                INNER JOIN users ON users.id = playlists.user_id
                INNER JOIN friendships ON friendships.friend_id = playlists.user_id
                WHERE friendships.user_id = :user_id
                AND friendships.pending = false
            ) AS t
            GROUP BY t.id, t.name, t.duration, t.spotify_uri, t.spotify_url, t.artists_json, t.album_json
            ORDER BY added_at DESC`

            // WHERE ` + strings.Join(filters, " AND ") + `

    return q
}

// Search - gets row(s) in tracks based on search params
func (ts *TrackStore) Search(params TracksSearchParams) ([]model.Track, error) {
    logger := logger.NewLogger()

    q := params.String()

    var tracks []model.Track
    nstmt, err := ts.DB.PrepareNamed(q)
    if err != nil {
        logger.Error().Err(err).Msg("TrackStore.Search - failed to prepare named statement")
        return tracks, err
    }
    err = nstmt.Select(&tracks, params)
    if err != nil {
        logger.Error().Err(err).Msg("TrackStore.Search - failed to get tracks")
        return tracks, err
    }

    return tracks, nil
}

// GetBySpotifyIDIfExists - get a row from tracks by Spotify ID if it exists
func (ts *TrackStore) GetBySpotifyIDIfExists(id string) (*model.Track, error) {
    logger := logger.NewLogger()

    q := `SELECT EXISTS (SELECT 1 FROM tracks WHERE spotify_id = $1)`

    var exists bool
    if err := ts.DB.QueryRowx(q, id).Scan(&exists); err != nil {
        logger.Error().
            Err(err).
            Str("spotifyID", id).
            Msg("TrackStore.GetBySpotifyIDIfExists - failed to check if track exists")
        return nil, err
    }
    if !exists {
        return nil, nil
    }

    q = `SELECT * FROM tracks WHERE spotify_id = $1`

    track := &model.Track{}
    if err := ts.DB.QueryRowx(q, id).StructScan(track); err != nil {
        logger.Error().
            Err(err).
            Str("spotifyID", id).
            Msg("TrackStore.GetBySpotifyIDIfExists - failed to get track")
        return nil, err
    }

    return track, nil
}

// Create - insert a row into tracks
func (ts *TrackStore) Create(t model.Track) (*model.Track, error) {
    logger := logger.NewLogger()

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
        logger.Error().
            Err(err).
            Str("trackID", t.ID).
            Msg("TrackStore.Create - failed to create track")
		return nil, err
	}

	return track, nil
}
