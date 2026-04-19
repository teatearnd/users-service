package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var secretKey string

func Init(secret string) error {
	s := strings.TrimSpace(secret)
	if s == "" {
		return fmt.Errorf("JWT secret is empty")
	}
	secretKey = s
	return nil
}

func ValidateConfig() error {
	if strings.TrimSpace(secretKey) == "" {
		return fmt.Errorf("JWT secret is not initialized")
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
