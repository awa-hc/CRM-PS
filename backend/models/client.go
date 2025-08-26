package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Email       string         `json:"email" gorm:"uniqueIndex"`
	Phone       string         `json:"phone"`
	Address     string         `json:"address"`
	City        string         `json:"city"`
	State       string         `json:"state"`
	ZipCode     string         `json:"zip_code"`
	Company     string         `json:"company"`
	TaxID       string         `json:"tax_id"` // RFC o identificación fiscal
	ContactType string         `json:"contact_type" gorm:"default:'individual'"` // individual, company
	Notes       string         `json:"notes" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Projects []Project `json:"projects,omitempty" gorm:"foreignKey:ClientID"`
	Quotes   []Quote   `json:"quotes,omitempty" gorm:"foreignKey:ClientID"`
	Invoices []Invoice `json:"invoices,omitempty" gorm:"foreignKey:ClientID"`
}

type CreateClientRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zip_code"`
	Company     string `json:"company"`
	TaxID       string `json:"tax_id"`
	ContactType string `json:"contact_type"`
	Notes       string `json:"notes"`
}

type UpdateClientRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"omitempty,email"`
	Phone       *string `json:"phone"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	ZipCode     *string `json:"zip_code"`
	Company     *string `json:"company"`
	TaxID       *string `json:"tax_id"`
	ContactType *string `json:"contact_type"`
	Notes       *string `json:"notes"`
	IsActive    *bool   `json:"is_active"`
}