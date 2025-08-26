package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Code         string         `json:"code" gorm:"uniqueIndex;not null"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description" gorm:"type:text"`
	ClientID     uint           `json:"client_id" gorm:"not null"`
	Client       Client         `json:"client" gorm:"foreignKey:ClientID"`
	Status       string         `json:"status" gorm:"default:'planning'"` // planning, in_progress, completed, cancelled, on_hold
	Priority     string         `json:"priority" gorm:"default:'medium'"` // low, medium, high, urgent
	Type         string         `json:"type" gorm:"default:'construction'"` // construction, renovation, maintenance
	ProjectType  string         `json:"project_type" gorm:"default:'construction'"` // construction, renovation, maintenance
	Address      string         `json:"address"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	ZipCode      string         `json:"zip_code"`
	StartDate    *time.Time     `json:"start_date"`
	EndDate      *time.Time     `json:"end_date"`
	Budget       float64        `json:"budget" gorm:"type:decimal(15,2)"`
	EstimatedCost float64       `json:"estimated_cost" gorm:"type:decimal(15,2);default:0"`
	ActualCost   float64        `json:"actual_cost" gorm:"type:decimal(15,2);default:0"`
	Progress     int            `json:"progress" gorm:"default:0"` // Porcentaje de 0-100
	Notes        string         `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Quotes          []Quote           `json:"quotes,omitempty" gorm:"foreignKey:ProjectID"`
	Invoices        []Invoice         `json:"invoices,omitempty" gorm:"foreignKey:ProjectID"`
	ProjectMaterials []ProjectMaterial `json:"project_materials,omitempty" gorm:"foreignKey:ProjectID"`
	WorkLogs        []WorkLog         `json:"work_logs,omitempty" gorm:"foreignKey:ProjectID"`
}

type CreateProjectRequest struct {
	Name          string     `json:"name" binding:"required"`
	Description   string     `json:"description"`
	ClientID      uint       `json:"client_id" binding:"required"`
	Status        string     `json:"status"`
	Priority      string     `json:"priority"`
	Type          string     `json:"type"`
	ProjectType   string     `json:"project_type"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	ZipCode       string     `json:"zip_code"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Budget        float64    `json:"budget"`
	EstimatedCost float64    `json:"estimated_cost"`
	Notes         string     `json:"notes"`
}

type UpdateProjectRequest struct {
	Name          *string    `json:"name"`
	Description   *string    `json:"description"`
	Status        *string    `json:"status"`
	Priority      *string    `json:"priority"`
	Type          *string    `json:"type"`
	ProjectType   *string    `json:"project_type"`
	Address       *string    `json:"address"`
	City          *string    `json:"city"`
	State         *string    `json:"state"`
	ZipCode       *string    `json:"zip_code"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Budget        *float64   `json:"budget"`
	EstimatedCost *float64   `json:"estimated_cost"`
	ActualCost    *float64   `json:"actual_cost"`
	Progress      *int       `json:"progress"`
	Notes         *string    `json:"notes"`
}

type ProjectStats struct {
	TotalProjects     int64   `json:"total_projects"`
	ActiveProjects    int64   `json:"active_projects"`
	CompletedProjects int64   `json:"completed_projects"`
	PlanningProjects  int64   `json:"planning_projects"`
	TotalBudget       float64 `json:"total_budget"`
	TotalActualCost   float64 `json:"total_actual_cost"`
	AverageProgress   float64 `json:"average_progress"`
}