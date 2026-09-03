package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Schedule string

const (
	ScheduleOTC Schedule = "OTC"
	ScheduleH   Schedule = "H" // requires prescription upload at checkout
)

// Product carries the homeopathy-specific fields called out in the brief:
// potency, form/dosage type, manufacturer, therapeutic category, indications,
// contraindications, and schedule (which drives the checkout flow).
type Product struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name              string         `gorm:"not null" json:"name"`               // e.g. Arnica Montana
	Slug              string         `gorm:"uniqueIndex;not null" json:"slug"`
	Potency           string         `json:"potency"`                            // 6C, 30C, 200C, 1M, Q, etc.
	Form              string         `json:"form"`                               // pellets/liquid/tablets/cream
	SizeQuantity      string         `json:"size_quantity"`                      // 10g, 30ml, etc.
	Manufacturer      string         `json:"manufacturer"`                       // SBL, Schwabe, Boiron...
	TherapeuticCategory string       `json:"therapeutic_category"`
	Indications       pq.StringArray       `gorm:"type:text[]" json:"indications"`
	Contraindications pq.StringArray       `gorm:"type:text[]" json:"contraindications"`
	Schedule          Schedule       `gorm:"type:varchar(5);default:'OTC'" json:"schedule"`
	HSNCode           string         `json:"hsn_code"`
	SKU               string         `gorm:"uniqueIndex;not null" json:"sku"`
	Price             int64          `gorm:"not null;check: price>=0" json:"price"` // stored in paise/cents to avoid float issues
	MRP               int64          `gorm:"check: mrp>=0" json:"mrp"`
	StockQty          int            `gorm:"default:0;check: stock_qty>=0" json:"stock_qty"`
	Images            pq.StringArray       `gorm:"type:text[]" json:"images"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
