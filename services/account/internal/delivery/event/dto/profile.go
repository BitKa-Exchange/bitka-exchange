package dto

import (
	"github.com/google/uuid"
)

type UserRegisteredEvent struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}
