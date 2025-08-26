package main

import (
	"log"
	"os"

	"raborimet-crm/backend/config"
	"raborimet-crm/backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// No cargar archivos .env - usar solo variables del sistema

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
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:4200",                // Angular dev server
		"https://crm-ps.vercel.app/auth/login", // Frontend en Vercel
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

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
