package models

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	InvoiceNumber string         `json:"invoice_number" gorm:"uniqueIndex;not null"`
	ClientID      uint           `json:"client_id" gorm:"not null"`
	Client        Client         `json:"client" gorm:"foreignKey:ClientID"`
	ProjectID     *uint          `json:"project_id"`
	Project       *Project       `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	QuoteID       *uint          `json:"quote_id"`
	Quote         *Quote         `json:"quote,omitempty" gorm:"foreignKey:QuoteID"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Status        string         `json:"status" gorm:"default:'draft'"` // draft, sent, paid, overdue, cancelled
	IssueDate     time.Time      `json:"issue_date" gorm:"not null"`
	DueDate       time.Time      `json:"due_date" gorm:"not null"`
	PaidDate      *time.Time     `json:"paid_date"`
	Subtotal      float64        `json:"subtotal" gorm:"type:decimal(15,2);default:0"`
	TaxRate       float64        `json:"tax_rate" gorm:"type:decimal(5,2);default:0"`
	TaxAmount     float64        `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	Discount      float64        `json:"discount" gorm:"type:decimal(15,2);default:0"`
	Total         float64        `json:"total" gorm:"type:decimal(15,2);default:0"`
	PaidAmount    float64        `json:"paid_amount" gorm:"type:decimal(15,2);default:0"`
	Balance       float64        `json:"balance" gorm:"type:decimal(15,2);default:0"`
	Notes         string         `json:"notes" gorm:"type:text"`
	Terms         string         `json:"terms" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Items []InvoiceItem `json:"items,omitempty" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`
}

type InvoiceItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	InvoiceID   uint    `json:"invoice_id" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`
	Unit        string  `json:"unit" gorm:"default:'pcs'"` // pcs, m2, m3, kg, etc.
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	Total       float64 `json:"total" gorm:"type:decimal(15,2);not null"`
	Notes       string  `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInvoiceRequest struct {
	ClientID    uint                     `json:"client_id" binding:"required"`
	ProjectID   *uint                    `json:"project_id"`
	QuoteID     *uint                    `json:"quote_id"`
	Title       string                   `json:"title" binding:"required"`
	Description string                   `json:"description"`
	DueDate     time.Time                `json:"due_date" binding:"required"`
	TaxRate     float64                  `json:"tax_rate"`
	Discount    float64                  `json:"discount"`
	Notes       string                   `json:"notes"`
	Terms       string                   `json:"terms"`
	Items       []CreateInvoiceItemRequest `json:"items" binding:"required,min=1"`
}

type CreateInvoiceItemRequest struct {
	Description string  `json:"description" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Unit        string  `json:"unit"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gte=0"`
	Notes       string  `json:"notes"`
}

type UpdateInvoiceRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	PaidDate    *time.Time `json:"paid_date"`
	TaxRate     float64    `json:"tax_rate"`
	Discount    float64    `json:"discount"`
	PaidAmount  float64    `json:"paid_amount"`
	Notes       string     `json:"notes"`
	Terms       string     `json:"terms"`
}

type InvoiceStats struct {
	TotalInvoices   int64   `json:"total_invoices"`
	DraftInvoices   int64   `json:"draft_invoices"`
	SentInvoices    int64   `json:"sent_invoices"`
	PaidInvoices    int64   `json:"paid_invoices"`
	OverdueInvoices int64   `json:"overdue_invoices"`
	TotalValue      float64 `json:"total_value"`
	PaidValue       float64 `json:"paid_value"`
	OutstandingValue float64 `json:"outstanding_value"`
}