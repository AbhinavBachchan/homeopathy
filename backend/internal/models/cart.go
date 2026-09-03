package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	GuestToken string    `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// One cart has many cart items.
	Items []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items"`
}
type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_cart_product" json:"cart_id"`
	Cart      Cart      `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"-"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_cart_product" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product"`
	Quantity  int       `gorm:"not null;check:quantity > 0" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
