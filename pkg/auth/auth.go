package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserCredentials struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var secretKey = os.Getenv("JWT_SECRET")

func ValidateConfig() error {
	if strings.TrimSpace(secretKey) == "" {
		return fmt.Errorf("JWT_SECRET is not configured in the .env")
	}
	return nil
}

// This function takes a secretKey from an .env "JWT_SECRET"
func CreateToken(email string) (string, *AccessClaims, error) {
	claims := AccessClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "users-service",
			Audience:  jwt.ClaimStrings{"surveys-service"}, // ?
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("failed when signing a token: %w", err)
	}
	return tokenString, &claims, nil
}

// Checks the issuer, the audience and validity of the token
func ValidateToken(tokenString string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString,
		claims,
		func(t *jwt.Token) (any, error) { return []byte(secretKey), nil },
		jwt.WithIssuer("users-service"),
		jwt.WithAudience("surveys-service"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return claims, nil
}
