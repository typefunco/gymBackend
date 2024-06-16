package models

import (
	"context"
	"fmt"
	"gymBackend/utils"

	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
)

type Trainer struct {
	Id               int `json:"Id"`
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

func GetTrainers() ([]Trainer, error) {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	SQLquery := "SELECT id, first_name, last_name, weight, height, fat_percentage, muscle_percentage, experience, education, schedule, gym_id FROM trainers"

	rows, err := conn.Query(context.Background(), SQLquery)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %v", err)
	}
	defer rows.Close()

	var trainers []Trainer
	for rows.Next() { // Going on the data. Id=1, Id=2, Id=3. Checking this data. If we have next element, we move, else stop
		var trainer Trainer
		var firstName, lastName, education, schedule pgtype.Varchar
		var weight, height, fatPercentage, musclePercentage pgtype.Float8
		var experience, gymID pgtype.Int4

		err := rows.Scan(&trainer.Id, &firstName, &lastName, &weight, &height, &fatPercentage, &musclePercentage, &experience, &education, &schedule, &gymID) // &trainer.Id because in DB it's int and here it's int
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if firstName.Status == pgtype.Present { // string
			trainer.FirstName = firstName.String
		}

		if lastName.Status == pgtype.Present { // string
			trainer.LastName = lastName.String
		}

		if weight.Status == pgtype.Present { // Float64
			trainer.Weight = weight.Float
		}

		if height.Status == pgtype.Present { // Float64
			trainer.Height = height.Float
		}

		if fatPercentage.Status == pgtype.Present { // Float64
			trainer.FatPercentage = fatPercentage.Float
		}

		if musclePercentage.Status == pgtype.Present { // Float64
			trainer.MusclePercentage = musclePercentage.Float
		}

		if experience.Status == pgtype.Present { // Int
			trainer.Experience = int(experience.Int)
		}

		if education.Status == pgtype.Present { // string
			trainer.Education = education.String
		}

		if schedule.Status == pgtype.Present { // string
			trainer.Schedule = schedule.String
		}

		if gymID.Status == pgtype.Present { // int
			trainer.GymId = int(gymID.Int)
		}

		trainers = append(trainers, trainer)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return trainers, nil
}

func (t *Trainer) UpdateProfile(trainerId int, updates map[string]interface{}) error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	sqlQuery := "UPDATE trainers SET"
	params := []interface{}{}
	paramId := 1

	for key, value := range updates {
		sqlQuery += fmt.Sprintf(" %s = $%d,", key, paramId)
		params = append(params, value)
		paramId++
	}
	sqlQuery = sqlQuery[:len(sqlQuery)-1]               // Remove the trailing comma
	sqlQuery += fmt.Sprintf(" WHERE id = $%d", paramId) // Use $1 for trainerId
	params = append(params, trainerId)

	_, err = conn.Exec(context.Background(), sqlQuery, params...)
	if err != nil {
		return fmt.Errorf("failed to update trainer profile: %v", err)
	}

	return nil
}
