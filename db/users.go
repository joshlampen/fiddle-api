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

func (us *UserStore) GetByAuthID(id string) (*model.User, error) {
    q := `SELECT * FROM users WHERE auth_id = $1`

    user := &model.User{}
    if err := us.DB.QueryRowx(q, id).StructScan(user); err != nil {
        return nil, err
    }

    return user, nil
}

// GetBySpotifyIDIfExists - get a row from users by Spotify ID if it exists
func (us *UserStore) GetBySpotifyIDIfExists(id string) (*model.User, error) {
    q := `SELECT EXISTS (SELECT 1 FROM users WHERE spotify_id = $1)`

    var exists bool
    if err := us.DB.QueryRowx(q, id).Scan(&exists); err != nil {
        return nil, err
    }
    if !exists {
        return nil, nil
    }

    q = `SELECT * FROM users WHERE spotify_id = $1`

    user := &model.User{}
    if err := us.DB.QueryRowx(q, id).StructScan(user); err != nil {
        return nil, err
    }

    return user, nil
}

// Create - insert a row into users
func (us *UserStore) Create(u model.User) (*model.User, error) {
	q := `INSERT INTO users (
			display_name,
			email,
			spotify_url,
			spotify_image_url,
			spotify_id,
            auth_id,
            token
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *`

    user := &model.User{}
	if err := us.DB.QueryRowx(
		q,
		u.DisplayName,
		u.Email,
		u.SpotifyURL,
		u.SpotifyImageURL,
		u.SpotifyID,
        u.AuthID,
        u.Token,
	).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}
