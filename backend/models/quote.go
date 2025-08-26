package models

import (
	"time"

	"gorm.io/gorm"
)

type Quote struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	QuoteNumber string         `json:"quote_number" gorm:"uniqueIndex;not null"`
	ClientID    uint           `json:"client_id" gorm:"not null"`
	Client      Client         `json:"client" gorm:"foreignKey:ClientID"`
	ProjectID   *uint          `json:"project_id"`
	Project     *Project       `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Status      string         `json:"status" gorm:"default:'draft'"` // draft, sent, accepted, rejected, expired
	ValidUntil  *time.Time     `json:"valid_until"`
	Subtotal    float64        `json:"subtotal" gorm:"type:decimal(15,2);default:0"`
	TaxRate     float64        `json:"tax_rate" gorm:"type:decimal(5,2);default:0"`
	TaxAmount   float64        `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	Discount    float64        `json:"discount" gorm:"type:decimal(15,2);default:0"`
	Total       float64        `json:"total" gorm:"type:decimal(15,2);default:0"`
	Notes       string         `json:"notes" gorm:"type:text"`
	Terms       string         `json:"terms" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Items []QuoteItem `json:"items,omitempty" gorm:"foreignKey:QuoteID;constraint:OnDelete:CASCADE"`
}

type QuoteItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	QuoteID     uint    `json:"quote_id" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`
	Unit        string  `json:"unit" gorm:"default:'pcs'"` // pcs, m2, m3, kg, etc.
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	Total       float64 `json:"total" gorm:"type:decimal(15,2);not null"`
	Notes       string  `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateQuoteRequest struct {
	ClientID    uint                    `json:"client_id" binding:"required"`
	ProjectID   *uint                   `json:"project_id"`
	Title       string                  `json:"title" binding:"required"`
	Description string                  `json:"description"`
	ValidUntil  *time.Time              `json:"valid_until"`
	TaxRate     float64                 `json:"tax_rate"`
	Discount    float64                 `json:"discount"`
	Notes       string                  `json:"notes"`
	Terms       string                  `json:"terms"`
	Items       []CreateQuoteItemRequest `json:"items" binding:"required,min=1"`
}

type CreateQuoteItemRequest struct {
	Description string  `json:"description" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Unit        string  `json:"unit"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gte=0"`
	Notes       string  `json:"notes"`
}

type UpdateQuoteRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	ValidUntil  *time.Time `json:"valid_until"`
	TaxRate     *float64   `json:"tax_rate"`
	Discount    *float64   `json:"discount"`
	Notes       *string    `json:"notes"`
	Terms       *string    `json:"terms"`
}

type QuoteStats struct {
	TotalQuotes      int64   `json:"total_quotes"`
	DraftQuotes      int64   `json:"draft_quotes"`
	SentQuotes       int64   `json:"sent_quotes"`
	AcceptedQuotes   int64   `json:"accepted_quotes"`
	RejectedQuotes   int64   `json:"rejected_quotes"`
	TotalValue       float64 `json:"total_value"`
	AcceptedValue    float64 `json:"accepted_value"`
	ThisMonthQuotes  int64   `json:"this_month_quotes"`
}