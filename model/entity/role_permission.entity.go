package entity

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoleID       uuid.UUID `json:"role_id" gorm:"type:uuid"`
	PermissionID uuid.UUID `json:"permission_id" gorm:"type:uuid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
