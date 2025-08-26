package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/routes"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables del sistema")
	}

	// Inicializar base de datos
	config.InitDB()

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Deshabilitar redirecciones automáticas de trailing slash
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"} // Angular dev server
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Configurar rutas
	routes.SetupRoutes(router)

	// Obtener puerto del entorno o usar 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en puerto %s", port)
	router.Run(":" + port)
}