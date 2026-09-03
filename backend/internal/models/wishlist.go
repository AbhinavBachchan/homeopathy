package models

import (
	"time"

	"github.com/google/uuid"
)

type WishlistItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_product" json:"user_id"`

	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_product" json:"product_id"`

	CreatedAt time.Time `json:"created_at"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}