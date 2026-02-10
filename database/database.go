package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	// "time"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	// Opsi A
	// db.SetConnMaxLifetime(30 * time.Minute)
	// db.SetConnMaxIdleTime(5 * time.Minute)
	// db.SetMaxOpenConns(25)

	// Opsi B
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	log.Println("Database connected successfully")
	return db, nil
}
