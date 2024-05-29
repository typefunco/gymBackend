package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// func InitDB() {
// 	var err error
// 	DB, err = sql.Open("sqlite3", "api.db")

// 	if err != nil {
// 		panic("Could not connect to the database.")
// 	}
// 	DB.SetMaxOpenConns(10)
// 	DB.SetMaxIdleConns(5)
// 	createTables()
// 	createAuthorDB()
// 	createUsersDB()
// }

// func createTables() {
// 	createEventsTable := `CREATE TABLE IF NOT EXISTS events (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		TotalPeople INTEGER,
// 		Theme TEXT NOT NULL,
// 		MinuteDuration INTEGER);
// 	`

// 	_, err := DB.Exec(createEventsTable)

// 	if err != nil {
// 		panic("Could not create events table.")
// 	}

// }
