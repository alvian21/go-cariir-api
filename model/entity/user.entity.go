package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FullName  string         `json:"full_name"`
	Email     string         `json:"email" gorm:"unique;not null"`
	IsActive  bool           `json:"is_active" gorm:"not null;default:true"`
	Password  string         `json:"-" gorm:"not null"`
	RoleID    *uuid.UUID     `json:"role_id" gorm:"type:uuid"`
	Role      *Role          `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}
