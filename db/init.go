package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/JoshLampen/fiddle/api/internal/constant"
)

// Initialize the database connection
func Init(port string) (*sqlx.DB, error) {
	var connString string
	if port == "8000" {
		connString = os.Getenv(constant.EnvVarDBConnectionInfo) //local db
	} else {
		connString = os.Getenv(constant.EnvVarDatabaseURL) // heroku db
	}
		
	// connString := os.Getenv(constant.EnvVarDBConnectionInfo) //local db
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("db.Init - Failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Init - Failed to connect to database: %w", err)
	}

	return db, nil
}
