package models

import (
	"context"
	"errors"
	"fmt"
	"gymBackend/utils"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id               int
	Username         string `json:"username"`
	Password         string `json:"password"`
	ProfileCreated   bool
	FirstName        string
	LastName         string
	Weight           float64
	Height           float64
	FatPercentage    float64
	MusclePercentage float64
}

var users []User

func (u *User) Save() error {
	dbURL, _ := utils.CheckDbConnection()

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
		INSERT INTO users (username, password, profile_created)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	// Set ProfileCreated to true in the SQL query
	err = conn.QueryRow(context.Background(), sqlQuery, u.Username, hashedPassword, true).Scan(&u.Id)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	u.ProfileCreated = true

	return nil
}

func (u *User) ValidateCredentials() error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	SQLquery := "SELECT id, password FROM users WHERE username = $1"
	row := conn.QueryRow(context.Background(), SQLquery, u.Username)

	var retrievePassword string

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

func GetUsers() ([]User, error) {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	SQLquery := "SELECT * FROM users"

	rows, err := conn.Query(context.Background(), SQLquery)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return users, nil

}

// func (u *User) SaveExtra() error {
// 	dbURL, _ := utils.CheckDbConnection()

// 	conn, err := pgx.Connect(context.Background(), dbURL)
// 	if err != nil {
// 		return fmt.Errorf("unable to connect to database: %v", err)
// 	}
// 	defer conn.Close(context.Background())

// 	sqlQuery := `
// 		INSERT INTO users (first_name, last_name, weight, height, fat_percentage, muscle_percentage)
// 		VALUES ($1, $2, $3, $4, $5, $6)
// 		RETURNING id
// 	`
// 	err = conn.QueryRow(context.Background(), sqlQuery, u.FirstName, u.LastName, u.Weight, u.Height, u.FatPercentage, u.MusclePercentage).Scan(&u.Id)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert user: %v", err)
// 	}

// 	return nil
// }
