package models

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type User struct {
	Id               int
	FirstName        string
	LastName         string
	Weight           float64
	Height           float64
	FatPercentage    float64
	MusclePercentage float64
}

func (u *User) Save() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	sqlQuery := `
		INSERT INTO users (first_name, last_name, weight, height, fat_percentage, muscle_percentage)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = conn.QueryRow(context.Background(), sqlQuery, u.FirstName, u.LastName, u.Weight, u.Height, u.FatPercentage, u.MusclePercentage).Scan(&u.Id)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}
