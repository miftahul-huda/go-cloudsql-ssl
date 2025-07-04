package models

import (
	"database/sql"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func Migrate(db *sql.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100)
	);`)
}
