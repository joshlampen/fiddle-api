package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/db/model"
	"github.com/JoshLampen/fiddle/api/internal/utils/logger"
)

// TokenStore manages the auth_codes database entity
type TokenStore struct {
	DB *sqlx.DB
}

// NewTokenStore returns a new db connection for the TokenStore
func NewTokenStore(db *sqlx.DB) *TokenStore {
	return &TokenStore{DB: db}
}

// GetByID - gets a row in auth_codes by ID
func (ts *TokenStore) GetByID(id string) (model.Token, error) {
    logger := logger.NewLogger()

    q := `SELECT * FROM tokens WHERE id = $1`

	var token model.Token
	if err := ts.DB.QueryRowx(q, id).StructScan(&token); err != nil {
        logger.Error().
            Err(err).
            Str("id", id).
            Msg("TokenStore.GetByID - failed to get token")
		return model.Token{}, err
	}

	return token, nil
}

// Create - insert a row into auth_codes
func (ts *TokenStore) Create(t model.Token) (model.Token, error) {
    logger := logger.NewLogger()

	q := `INSERT INTO tokens (
			id,
			access_token,
            token_type,
            scope,
            expires_in,
            refresh_token
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *`

	var token model.Token
	if err := ts.DB.QueryRowx(
		q,
		t.ID,
		t.AccessToken,
        t.TokenType,
        t.Scope,
        t.ExpiresIn,
        t.RefreshToken,
	).StructScan(&token); err != nil {
        logger.Error().
            Err(err).
            Str("id", t.ID).
            Msg("TokenStore.Create - failed to create token")
		return model.Token{}, err
	}

	return token, nil
}
