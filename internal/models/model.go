package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `json:"user_id" db:"id"`
	FirstName sql.NullString `json:"first_name" db:"first_name"`
	LastName  sql.NullString `json:"last_name" db:"last_name"`
	Login     string         `json:"login" db:"login"`
	Password  string         `json:"-" db:"password"`
	LastLogin time.Time      `json:"last_login" db:"last_login"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
}
type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty" binding:"omitempty,min=2"`
	LastName  string `json:"last_name,omitempty" binding:"omitempty,min=2"`
	Login     string `json:"login,omitempty" binding:"omitempty,min=3"`

	Password string `json:"password,omitempty" binding:"omitempty,min=8"`
}
type ServerResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
