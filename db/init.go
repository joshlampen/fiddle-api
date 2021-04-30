package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/JoshLampen/fiddle/api/internal/constant"
)

// Initialize the database connection
func Init(port string) (*sqlx.DB, error) {
	var connString string
	if port == "8001" {
		if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
			return nil, err
		}
		connString = os.Getenv(constant.EnvVarDBConnectionInfo) //local db
	} else {
		connString = os.Getenv(constant.EnvVarDatabaseURL) // heroku db
	}
		
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
