package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"raborimet-crm/backend/models"
)

var DB *gorm.DB

func InitDB() {
	// Configuración de la base de datos
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "raborimet_crm")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Cadena de conexión
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

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