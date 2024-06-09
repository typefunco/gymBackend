package models

import (
	"context"
	"fmt"
	"gymBackend/utils"

	"github.com/jackc/pgx/v5"
)

type GymInfo struct {
	Id           int64
	Address      string
	NearestMetro string
	NumMachine   int64
	NumTrainers  int64
	NumUsers     int64
	Area         float64
}

func (gi *GymInfo) Save() error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	sqlQuery := `
		INSERT INTO users (address, nearest_metro, area)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = conn.QueryRow(context.Background(), sqlQuery, gi.Address, gi.NearestMetro, gi.Area).Scan(&gi.Id)
	if err != nil {
		return fmt.Errorf("FAILED TO SAVE GYM INFO: %v", err)
	}

	return nil
}
