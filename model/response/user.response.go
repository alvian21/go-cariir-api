package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	FullName  string     `json:"full_name"`
	Email     string     `json:"email"`
	IsActive  bool       `json:"is_active"`
	RoleID    *uuid.UUID `json:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
