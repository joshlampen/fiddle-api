package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/JoshLampen/fiddle/api/internal/constant"
	"github.com/JoshLampen/fiddle/api/internal/utils/logger"
)

// Initialize the database connection
func Init(port string) (*sqlx.DB, error) {
    logger := logger.NewLogger()

	var connString string
	if port == "8001" {
		if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
            logger.Error().Err(err).Msg("db.Init - failed to load local env file")
			return nil, err
		}
		connString = os.Getenv(constant.EnvVarDBConnectionInfo) //local db
	} else {
		connString = os.Getenv(constant.EnvVarDatabaseURL) // heroku db
	}

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
        logger.Error().Err(err).Msg("db.Init - failed to open postgres database")
		return nil, err
	}

	if err := db.Ping(); err != nil {
        logger.Error().Err(err).Msg("db.Init - failed to connect to postgres database")
		return nil, err
	}

	return db, nil
}
