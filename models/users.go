package models

import (
	"context"
	"errors"
	"fmt"
	"gymBackend/utils"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type User struct {
	Id               int
	Username         string `json:"username"`
	Password         string `json:"password"`
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
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("CAN'T HASH PASSWORD: %v", err)
	}

	sqlQuery := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`
	err = conn.QueryRow(context.Background(), sqlQuery, u.Username, hashedPassword).Scan(&u.Id)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

func (u *User) ValidateCredentials() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	SQLquery := "SELECT id, password_hash FROM users WHERE username = $1"
	row := conn.QueryRow(context.Background(), SQLquery, u.Username)

	// Declare variables to store retrieved data
	var retrievePassword string

	// Scan the retrieved data into variables
	err = row.Scan(&u.Id, &retrievePassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("USER NOT FOUND")
		}
		return errors.New("CREDENTIALS INVALID")
	}

	isValidPassword := utils.CheckPasswordHash(u.Password, retrievePassword)
	if !isValidPassword {
		return errors.New("CREDENTIALS INVALID")
	}

	return nil
}

func (u *User) SaveExtra() error {
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
