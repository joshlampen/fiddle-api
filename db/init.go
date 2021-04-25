package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/JoshLampen/fiddle/api/internal/constant"
)

// Initialize the database connection
func Init() (*sqlx.DB, error) {
	if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
		return nil, fmt.Errorf("db.Init - Failed to load .env file: %w", err)
	}
	psqlInfo := os.Getenv(constant.EnvVarDBConnectionInfo)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("db.Init - Failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Init - Failed to connect to database: %w", err)
	}

	return db, nil
}
