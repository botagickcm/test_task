package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `json:"user_id" db:"user_id" example:"1"`
	FirstName sql.NullString `json:"first_name" db:"first_name" swaggertype:"string" example:"Дина"`
	LastName  sql.NullString `json:"last_name" db:"last_name" swaggertype:"string" example:"Алмазова"`
	Login     string         `json:"login" db:"login" example:"dina@gmail.com"`
	Password  string         `json:"-" db:"password"`
	LastLogin time.Time      `json:"last_login" db:"last_login" example:"2026-01-28T12:30:45+05:00"`
	CreatedAt time.Time      `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" db:"updated_at"`
	IsActive  bool           `json:"is_active,omitempty" db:"is_active"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required" example:"dina@gmail.com"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required" example:"Дина"`
	LastName  string `json:"last_name" binding:"required" example:"Алмазова"`
	Login     string `json:"login" binding:"required,email" example:"dina@gmail.com"`
	Password  string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty" binding:"omitempty,min=2" example:"Дана"`
	LastName  string `json:"last_name,omitempty" binding:"omitempty,min=2" example:"Алмасова"`
	Login     string `json:"login,omitempty" binding:"omitempty,min=3" example:"dina1@gmail.com"`

	Password string `json:"password,omitempty" binding:"omitempty,min=8" example:"NewPassword456!"`
}

type ServerResponse struct {
	Status  string      `json:"status" `
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  *User  `json:"user"`
}
