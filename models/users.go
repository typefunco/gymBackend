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
	TrainerId        int
	IsSuperUser      bool
}

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

	SQLquery := "SELECT id, username, password FROM users"

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

func (u *User) UpdateProfile(userId int, updates map[string]interface{}) error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	sqlQuery := "UPDATE users SET"
	params := []interface{}{}
	paramId := 1

	for key, value := range updates {
		sqlQuery += fmt.Sprintf(" %s = $%d,", key, paramId) // making SQL query
		params = append(params, value)
		paramId++
	}
	sqlQuery = sqlQuery[:len(sqlQuery)-1] // Remove the trailing comma
	sqlQuery += fmt.Sprintf(" WHERE id = $%d", paramId)
	params = append(params, userId)

	_, err = conn.Exec(context.Background(), sqlQuery, params...) // dynamic unpack
	if err != nil {
		return fmt.Errorf("failed to update user profile: %v", err)
	}

	return nil
}
