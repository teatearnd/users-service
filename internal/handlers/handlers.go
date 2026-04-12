package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/internal/dto"
	"example.com/internal/validations"
	"example.com/pkg/auth"
)

type Handler struct {
	DB *sql.DB
}

// todo use methods h *handler & remove hardcoded credentials
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userCredentials auth.UserCredentials

	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		http.Error(w, "failed to decode a request", http.StatusBadRequest)
		return
	}
	if userCredentials.Email == "xyecoc@vilniustech.lt" && userCredentials.PasswordHash == "hash123" {
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
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid credentials")
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userCredentials dto.UserRegistration

	if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
		http.Error(w, "failed to decode a request", http.StatusBadRequest)
		return
	}
	err := validations.ValidateEmail(userCredentials.Email)
	if err != nil {
		http.Error(w, "failed to validate the email", http.StatusBadRequest)
		return
	}
	err = validations.ValidatePassword(userCredentials.Password)
	if err != nil {
		http.Error(w, "failed to validate the password", http.StatusBadRequest)
		return
	}
}
