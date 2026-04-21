package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"example.com/internal/dto"
	"example.com/internal/repository"
	"example.com/internal/validations"
	"example.com/pkg/auth"
)

type Handler struct {
	DB             *sql.DB
	AllowedDomains []string
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userCredentials dto.UserCredentials

	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		http.Error(w, "failed to decode a request", http.StatusBadRequest)
		return
	}
	hash, err := repository.FindUserByEmail(h.DB, userCredentials.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if !auth.CheckPassword(userCredentials.Password, hash) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, claims, err := auth.CreateToken(userCredentials.Email)
	if err != nil {
		http.Error(w, "failed to create a token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "logged in",
		"token":   tokenString,
		"email":   claims.Email,
		"expires": claims.ExpiresAt.Time,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode a response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userCredentials dto.UserRegistration

	if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
		http.Error(w, "failed to decode a request", http.StatusBadRequest)
		return
	}
	err := validations.ValidateEmail(userCredentials.Email, h.AllowedDomains)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validations.ValidatePassword(userCredentials.Password)
	if err != nil {
		http.Error(w, "failed to validate the password", http.StatusBadRequest)
		return
	}
	// hash
	hashed, err := auth.HashPassword(userCredentials.Password)
	if err != nil {
		http.Error(w, "failed to hash the password", http.StatusInternalServerError)
		return
	}
	userCredentials.Password = hashed // very poor way of hashing but what can you do

	err = repository.CreateUser(h.DB, userCredentials)
	if err != nil {
		http.Error(w, "failed to insert a user into the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
