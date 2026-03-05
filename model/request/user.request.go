package request

import "github.com/google/uuid"

type UserCreateRequest struct {
	FullName string     `json:"full_name" validate:"required"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=6"`
	RoleID   *uuid.UUID `json:"role_id"`
}

type UserUpdateRequest struct {
	FullName string `json:"full_name" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type UserEmailRequest struct {
	Email string `json:"email" validate:"required"`
}
