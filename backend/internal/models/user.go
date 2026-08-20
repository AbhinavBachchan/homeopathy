package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin       Role = "admin"
	RolePatient     Role = "patient"
	RoleDoctor      Role = "doctor"
	RoleCorporateHR Role = "corporate_hr"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	Phone         string    `gorm:"index" json:"phone"`
	Name          string    `json:"name"`
	Role          Role      `gorm:"type:varchar(20);not null;default:'patient'" json:"role"`
	GoogleOAuthID string    `json:"-"`
	IsVerified    bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
