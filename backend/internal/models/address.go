package models

import (
	"time"

	"github.com/google/uuid"
)

type AddressType string

const (
	AddressTypeHome  AddressType = "home"
	AddressTypeWork  AddressType = "work"
	AddressTypeOther AddressType = "other"
)

type Address struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	// Relationship with User
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	AddressType AddressType `gorm:"type:varchar(20);not null;default:'home'" json:"address_type"`

	FullName string `gorm:"not null" json:"full_name"`
	Phone    string `gorm:"not null" json:"phone"`

	AddressLine1 string `gorm:"not null" json:"address_line1"`
	AddressLine2 string `json:"address_line2"`

	City    string `gorm:"not null" json:"city"`
	State   string `gorm:"not null" json:"state"`
	Country string `gorm:"not null;default:'India'" json:"country"`
	Pincode string `gorm:"not null" json:"pincode"`

	IsDefault bool `gorm:"not null;default:false" json:"is_default"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
