package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	FullName  string     `json:"fullName"`
	Email     string     `json:"email"`
	IsActive  bool       `json:"isActive"`
	RoleID    *uuid.UUID `json:"roleId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
