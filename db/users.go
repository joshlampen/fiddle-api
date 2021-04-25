package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/db/model"
)

// UserStore manages the tracks database entity
type UserStore struct {
	DB *sqlx.DB
}

// NewUserStore returns a new db connection for the UserStore
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{DB: db}
}

// Create - insert a row into users
func (us *UserStore) Create(u model.User) (model.User, error) {
	q := `INSERT INTO users (
			display_name,
			email,
			spotify_url,
			spotify_image_url,
			spotify_id,
            token
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *`

    var user model.User
	if err := us.DB.QueryRowx(
		q,
		u.DisplayName,
		u.Email,
		u.SpotifyURL,
		u.SpotifyImageURL,
		u.SpotifyID,
        u.Token,
	).StructScan(&user); err != nil {
		return model.User{}, err
	}

	return user, nil
}
