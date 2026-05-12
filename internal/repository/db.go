package repository

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/internal/config"
	"example.com/internal/dto"

	_ "github.com/lib/pq"
)

func OpenDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseUrl)
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

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		return fmt.Errorf("failed to get pgcrypto")
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),	
		password_hash TEXT NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		return fmt.Errorf("failed to initialize users-table: %w", err)
	}
	_, err = db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';`)
	if err != nil {
		return fmt.Errorf("failed to ensure role column: %w", err)
	}
	return nil
}

func CreateUser(h *sql.DB, cred dto.UserRegistration) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`
	_, err := h.Exec(query, cred.Email, cred.Password)
	if err != nil {
		return fmt.Errorf("failed to insert a user: %w", err)
	}
	return nil
}

func FindUserByEmail(h *sql.DB, email string) (string, error) {
	query := `SELECT password_hash FROM users WHERE email = $1`
	var hash string
	err := h.QueryRow(query, email).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	return hash, nil
}

func FindUserForLogin(h *sql.DB, email string) (string, string, string, error) {
	query := `SELECT id, role, password_hash FROM users WHERE email = $1`
	var id string
	var role string
	var hash string
	err := h.QueryRow(query, email).Scan(&id, &role, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", fmt.Errorf("user not found")
		}
		return "", "", "", err
	}
	return id, role, hash, nil
}
