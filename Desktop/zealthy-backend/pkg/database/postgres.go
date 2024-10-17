package database

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func NewPostgresConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
