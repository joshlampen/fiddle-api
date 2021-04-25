package db

import "github.com/jmoiron/sqlx"

// Store is an interface containing the connections to all database entities
type Store struct {
    *TokenStore
	*UserStore
	*PlaylistStore
	*TrackStore
    *PlaylistTrackStore
}

// NewStore returns a new instance of the Store
func NewStore(db *sqlx.DB) *Store {
	return &Store{
        TokenStore: NewTokenStore(db),
		UserStore: NewUserStore(db),
		PlaylistStore: NewPlaylistStore(db),
		TrackStore: NewTrackStore(db),
        PlaylistTrackStore: NewPlaylistTrackStore(db),
	}
}
