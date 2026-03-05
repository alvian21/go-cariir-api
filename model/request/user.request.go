package request

import "github.com/google/uuid"

type UserCreateRequest struct {
	FullName string     `json:"fullName" validate:"required"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=6"`
	RoleID   *uuid.UUID `json:"roleId"`
}

type UserUpdateRequest struct {
	FullName string `json:"fullName" validate:"required"`
	IsActive bool   `json:"isActive"`
}

type UserEmailRequest struct {
	Email string `json:"email" validate:"required"`
}
