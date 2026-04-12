package repository

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/internal/dto"
	_ "github.com/lib/pq"
)

func OpenDB() (*sql.DB, error) {
	dsn := "host=localhost port=5432 user=users_app password=users_pass dbname=users_service sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the db: %w", err)
	}
	// maxopenconns, maxidleconns?
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}
	log.Printf("established connection to db")
	return db, nil
}

func Registration(h *sql.DB, cred dto.UserRegistration) {
	// todo
}
