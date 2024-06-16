package models

import (
	"context"
	"fmt"
	"gymBackend/utils"

	"github.com/jackc/pgx/v5"
)

type Trainer struct {
	Id               int
	FirstName        string
	LastName         string
	Weight           float64
	Height           float64
	FatPercentage    float64
	MusclePercentage float64
	Experience       int
	Education        string
	Schedule         string
	GymId            int
}

func (t *Trainer) Save() error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	sqlQuery := `
		INSERT INTO trainers (first_name, last_name, weight, height, fat_percentage, muscle_percentage, experience, education, schedule, gym_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err = conn.QueryRow(context.Background(), sqlQuery, t.FirstName, t.LastName, t.Weight, t.Height, t.FatPercentage, t.MusclePercentage, t.Experience, t.Education, t.Schedule, t.GymId).Scan(&t.Id)
	if err != nil {
		return fmt.Errorf("failed to save trainer: %v", err)
	}

	return nil
}
