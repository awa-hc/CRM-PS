package models

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Category    string         `json:"category" gorm:"not null"` // cement, steel, wood, electrical, plumbing, etc.
	Unit        string         `json:"unit" gorm:"not null"` // kg, m3, m2, pcs, etc.
	UnitPrice   float64        `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	Supplier    string         `json:"supplier"`
	SKU         string         `json:"sku" gorm:"uniqueIndex"`
	Stock       float64        `json:"stock" gorm:"type:decimal(10,2);default:0"`
	MinStock    float64        `json:"min_stock" gorm:"type:decimal(10,2);default:0"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	Notes       string         `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	ProjectMaterials []ProjectMaterial `json:"project_materials,omitempty" gorm:"foreignKey:MaterialID"`
}

type ProjectMaterial struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	ProjectID        uint      `json:"project_id" gorm:"not null"`
	Project          Project   `json:"project" gorm:"foreignKey:ProjectID"`
	MaterialID       uint      `json:"material_id" gorm:"not null"`
	Material         Material  `json:"material" gorm:"foreignKey:MaterialID"`
	QuantityPlanned  float64   `json:"quantity_planned" gorm:"type:decimal(10,2);not null"`
	QuantityUsed     float64   `json:"quantity_used" gorm:"type:decimal(10,2);default:0"`
	UnitPrice        float64   `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	TotalCost        float64   `json:"total_cost" gorm:"type:decimal(15,2);not null"`
	Status           string    `json:"status" gorm:"default:'planned'"` // planned, ordered, delivered, used
	DeliveryDate     *time.Time `json:"delivery_date"`
	Notes            string    `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type WorkLog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ProjectID   uint           `json:"project_id" gorm:"not null"`
	Project     Project        `json:"project" gorm:"foreignKey:ProjectID"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Date        time.Time      `json:"date" gorm:"not null"`
	StartTime   time.Time      `json:"start_time" gorm:"not null"`
	EndTime     time.Time      `json:"end_time" gorm:"not null"`
	Hours       float64        `json:"hours" gorm:"type:decimal(5,2);not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	WorkType    string         `json:"work_type" gorm:"not null"` // construction, planning, supervision, etc.
	HourlyRate  float64        `json:"hourly_rate" gorm:"type:decimal(10,2)"`
	TotalCost   float64        `json:"total_cost" gorm:"type:decimal(15,2)"`
	Notes       string         `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateMaterialRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category" binding:"required"`
	Unit        string  `json:"unit" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gte=0"`
	Supplier    string  `json:"supplier"`
	SKU         string  `json:"sku"`
	Stock       float64 `json:"stock"`
	MinStock    float64 `json:"min_stock"`
	Notes       string  `json:"notes"`
}

type CreateProjectMaterialRequest struct {
	MaterialID      uint       `json:"material_id" binding:"required"`
	QuantityPlanned float64    `json:"quantity_planned" binding:"required,gt=0"`
	UnitPrice       float64    `json:"unit_price" binding:"required,gte=0"`
	Status          string     `json:"status"`
	DeliveryDate    *time.Time `json:"delivery_date"`
	Notes           string     `json:"notes"`
}

type UpdateMaterialRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	Unit        *string  `json:"unit"`
	UnitPrice   *float64 `json:"unit_price"`
	Supplier    *string  `json:"supplier"`
	SKU         *string  `json:"sku"`
	Stock       *float64 `json:"stock"`
	MinStock    *float64 `json:"min_stock"`
	IsActive    *bool    `json:"is_active"`
	Notes       *string  `json:"notes"`
}

type CreateWorkLogRequest struct {
	ProjectID   uint      `json:"project_id" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Description string    `json:"description" binding:"required"`
	WorkType    string    `json:"work_type" binding:"required"`
	HourlyRate  float64   `json:"hourly_rate"`
	Notes       string    `json:"notes"`
}