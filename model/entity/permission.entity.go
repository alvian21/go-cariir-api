package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Alias     string    `json:"alias"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Roles     []Role    `json:"roles" gorm:"many2many:role_permissions;"`
}
