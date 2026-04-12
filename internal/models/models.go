package models

import (
	"github.com/google/uuid"
)

// useless remove
type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
