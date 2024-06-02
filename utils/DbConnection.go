package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/withmandala/go-log"
)

func CheckConnection() bool {
	logger := log.New(os.Stderr)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
		return false
	}
	defer conn.Close(context.Background())

	return true
}
