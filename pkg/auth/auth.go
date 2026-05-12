package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Settings struct {
	Secret   string
	Issuer   string
	Audience string
}

var secretKey string
var issuer string
var audience string

func Init(s Settings) error {
	secretKey = s.Secret
	issuer = s.Issuer
	audience = s.Audience
	return nil
}

func ValidateConfig() error {
	if strings.TrimSpace(secretKey) == "" {
		return fmt.Errorf("JWT secret is not initialized")
	}
	return nil
}

// This function takes a secretKey from an .env "JWT_SECRET"
func CreateToken(email string, userID string, role string) (string, *AccessClaims, error) {
	claims := AccessClaims{
		Email:  email,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
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
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return claims, nil
}
