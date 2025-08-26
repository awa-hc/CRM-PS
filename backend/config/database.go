package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"raborimet-crm/backend/models"
)

var DB *gorm.DB

func InitDB() {
	// Usar DATABASE_URL directamente o valor por defecto
	defaultDatabaseURL := "postgres://avnadmin:AVNS_6JAvyq33eXv1P0af8pJ@pg-35eb9e99-leohermoso18-d7f0.j.aivencloud.com:12238/defaultdb?sslmode=require"
	dsn := getEnv("DATABASE_URL", defaultDatabaseURL)

	// Conectar a la base de datos
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}

	log.Println("Conexión a la base de datos establecida")

	// Ejecutar migraciones automáticas
	runMigrations()
}

func runMigrations() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Project{},
		&models.Quote{},
		&models.QuoteItem{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.Material{},
		&models.ProjectMaterial{},
		&models.WorkLog{},
	)
	if err != nil {
		log.Fatal("Error en las migraciones:", err)
	}
	log.Println("Migraciones ejecutadas correctamente")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}