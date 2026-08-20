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

type OrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Potency   string    `json:"potency"`
	Qty       int       `json:"qty"`
	Price     int64     `json:"price"`
}

type Order struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID            uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Status            OrderStatus    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Items             []byte         `gorm:"type:jsonb" json:"-"` // serialized []OrderItem
	Subtotal          int64          `json:"subtotal"`
	GSTAmount         int64          `json:"gst_amount"`
	ShippingCharge    int64          `json:"shipping_charge"`
	Total             int64          `json:"total"`
	PaymentMethod     string         `json:"payment_method"`
	PaymentID         string         `json:"payment_id"`
	AddressID         uuid.UUID      `gorm:"type:uuid" json:"address_id"`
	ShiprocketOrderID string         `json:"shiprocket_order_id"`
	TrackingNumber    string         `json:"tracking_number"`
	PrescriptionURL   string         `json:"prescription_url"` // required if any item is Schedule H
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}
