package models

// "github.com/jackc/pgx/v5"

type GymInfo struct {
	ID           int64
	Address      string
	NearestMetro string
	NumMachine   int64
	NumTrainers  int64
	NumUsers     int64
	Area         float64
}
