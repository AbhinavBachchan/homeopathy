package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderConfirmed OrderStatus = "confirmed"
	OrderShipped   OrderStatus = "shipped"
	OrderDelivered OrderStatus = "delivered"
	OrderCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID     uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	User   User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Status OrderStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	// One order has many order items
	Items             []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
	Subtotal          int64       `gorm:"not null;check:subtotal >= 0" json:"subtotal"`
	GSTAmount         int64       `gorm:"not null;default:0;check:gst_amount >= 0" json:"gst_amount"`
	ShippingCharge    int64       `gorm:"not null;default:0;check:shipping_charge >= 0" json:"shipping_charge"`
	Total             int64       `gorm:"not null;check:total >= 0" json:"total"`
	PaymentMethod     string      `json:"payment_method"`
	PaymentID         string      `json:"payment_id"`
	AddressID         uuid.UUID   `gorm:"type:uuid;not null;index" json:"address_id"`
	Address           Address     `gorm:"foreignKey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"address"`
	ShiprocketOrderID string      `json:"shiprocket_order_id"`
	TrackingNumber    string      `json:"tracking_number"`
	PrescriptionURL   string      `json:"prescription_url"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	// Relationship with Order
	OrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Order   Order     `gorm:"foreignKey:OrderID" json:"-"`
	// Relationship with Product
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	// Product snapshot at the time of purchase
	Name      string    `gorm:"not null" json:"name"`
	Potency   string    `json:"potency"`
	Qty       int       `gorm:"not null;check:qty > 0" json:"qty"`
	Price     int64     `gorm:"not null;check:price >= 0" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
