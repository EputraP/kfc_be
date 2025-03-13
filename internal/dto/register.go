package dto

import "github.com/google/uuid"

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}
