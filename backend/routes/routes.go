package routes

import (
	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/controllers"
	"raborimet-crm/backend/middleware"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(router *gin.Engine) {
	// Inicializar controladores
	authController := controllers.NewAuthController()
	clientController := controllers.NewClientController()
	projectController := controllers.NewProjectController()
	quoteController := controllers.NewQuoteController()
	materialController := controllers.NewMaterialController()
	dashboardController := controllers.NewDashboardController()
	reportController := controllers.NewReportController()

	// Grupo de rutas de la API
	api := router.Group("/api/v1")
	{
		// Rutas públicas de autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authController.Logout)
		}

		// Rutas protegidas que requieren autenticación
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Rutas de autenticación protegidas
			protectedAuth := protected.Group("/auth")
			{
				protectedAuth.GET("/profile", authController.GetProfile)
				protectedAuth.PUT("/profile", authController.UpdateProfile)
				protectedAuth.POST("/change-password", authController.ChangePassword)
				protectedAuth.GET("/verify", authController.VerifyToken)
			}

			// Rutas de clientes
			clients := protected.Group("/clients")
			{
				clients.GET("", clientController.GetClients)
				clients.GET("/stats", clientController.GetClientStats)
				clients.GET("/:id", clientController.GetClient)
				clients.POST("", clientController.CreateClient)
				clients.PUT("/:id", clientController.UpdateClient)
				clients.DELETE("/:id", clientController.DeleteClient)
			}

			// Rutas de proyectos
			projects := protected.Group("/projects")
			{
				projects.GET("", projectController.GetProjects)
				projects.GET("/stats", projectController.GetProjectStats)
				projects.GET("/:id", projectController.GetProject)
				projects.GET("/:id/materials", projectController.GetProjectMaterials)
				projects.POST("", projectController.CreateProject)
				projects.PUT("/:id", projectController.UpdateProject)
				projects.DELETE("/:id", projectController.DeleteProject)
			}

			// Rutas de cotizaciones
			quotes := protected.Group("/quotes")
			{
				quotes.GET("", quoteController.GetQuotes)
				quotes.GET("/stats", quoteController.GetQuoteStats)
				quotes.GET("/:id", quoteController.GetQuote)
				quotes.POST("", quoteController.CreateQuote)
				quotes.PUT("/:id", quoteController.UpdateQuote)
				quotes.DELETE("/:id", quoteController.DeleteQuote)
				quotes.PATCH("/:id/status", quoteController.ChangeQuoteStatus)
			}

			// Rutas de facturas (placeholder para futuras implementaciones)
			invoices := protected.Group("/invoices")
			{
				// TODO: Implementar controlador de facturas
				invoices.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Invoices endpoint - Coming soon"})
				})
			}

			// Rutas de materiales
			materials := protected.Group("/materials")
			{
				materials.GET("", materialController.GetMaterials)
				materials.GET("/stats", materialController.GetMaterialStats)
				materials.GET("/low-stock", materialController.GetLowStockMaterials)
				materials.GET("/categories", materialController.GetMaterialCategories)
				materials.GET("/:id", materialController.GetMaterial)
				materials.POST("", materialController.CreateMaterial)
				materials.PUT("/:id", materialController.UpdateMaterial)
				materials.DELETE("/:id", materialController.DeleteMaterial)
				materials.PATCH("/:id/stock", materialController.UpdateMaterialStock)
			}

			// Rutas de reportes
			reports := protected.Group("/reports")
			{
				reports.GET("/clients", reportController.GetClientsReport)
				reports.GET("/projects", reportController.GetProjectsReport)
				reports.GET("/quotes", reportController.GetQuotesReport)
				reports.GET("/materials", reportController.GetMaterialsReport)
				reports.GET("/financial", reportController.GetFinancialReport)
			}

			// Rutas de dashboard
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", dashboardController.GetDashboardStats)
				dashboard.GET("/recent-activity", dashboardController.GetRecentActivity)
				dashboard.GET("/projects-by-status", dashboardController.GetProjectsByStatus)
				dashboard.GET("/monthly-revenue", dashboardController.GetMonthlyRevenue)
				dashboard.GET("/upcoming-deadlines", dashboardController.GetUpcomingDeadlines)
				dashboard.GET("/financial-summary", dashboardController.GetFinancialSummary)
			}

			// Rutas administrativas (solo para administradores)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				// Gestión de usuarios
				users := admin.Group("/users")
				{
					users.GET("/", getUsersList)
					users.GET("/:id", getUserByID)
					users.PUT("/:id", updateUser)
					users.DELETE("/:id", deleteUser)
					users.POST("/:id/activate", activateUser)
					users.POST("/:id/deactivate", deactivateUser)
				}

				// Configuración del sistema
				system := admin.Group("/system")
				{
					system.GET("/info", getSystemInfo)
					system.GET("/health", getHealthCheck)
				}
			}
		}
	}

	// Ruta de health check pública
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "Raborimet CRM API",
			"version": "1.0.0",
		})
	})

	// Ruta de información de la API
	router.GET("/api/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Raborimet CRM API",
			"version":     "1.0.0",
			"description": "API para el sistema CRM de construcción Raborimet",
			"endpoints": gin.H{
				"auth":      "/api/v1/auth",
				"clients":   "/api/v1/clients",
				"projects":  "/api/v1/projects",
				"quotes":    "/api/v1/quotes",
				"invoices":  "/api/v1/invoices",
				"materials": "/api/v1/materials",
				"reports":   "/api/v1/reports",
				"dashboard": "/api/v1/dashboard",
			},
		})
	})
}

// Handlers para funcionalidades administrativas

// getUsersList obtiene lista de usuarios (solo admin)
func getUsersList(c *gin.Context) {
	// TODO: Implementar lógica de lista de usuarios
	c.JSON(200, gin.H{
		"message": "Users list - Coming soon",
		"users":   []interface{}{},
	})
}

// getUserByID obtiene usuario por ID (solo admin)
func getUserByID(c *gin.Context) {
	// TODO: Implementar lógica de obtener usuario por ID
	c.JSON(200, gin.H{
		"message": "Get user by ID - Coming soon",
	})
}

// updateUser actualiza usuario (solo admin)
func updateUser(c *gin.Context) {
	// TODO: Implementar lógica de actualizar usuario
	c.JSON(200, gin.H{
		"message": "Update user - Coming soon",
	})
}

// deleteUser elimina usuario (solo admin)
func deleteUser(c *gin.Context) {
	// TODO: Implementar lógica de eliminar usuario
	c.JSON(200, gin.H{
		"message": "Delete user - Coming soon",
	})
}

// activateUser activa usuario (solo admin)
func activateUser(c *gin.Context) {
	// TODO: Implementar lógica de activar usuario
	c.JSON(200, gin.H{
		"message": "Activate user - Coming soon",
	})
}

// deactivateUser desactiva usuario (solo admin)
func deactivateUser(c *gin.Context) {
	// TODO: Implementar lógica de desactivar usuario
	c.JSON(200, gin.H{
		"message": "Deactivate user - Coming soon",
	})
}

// getSystemInfo obtiene información del sistema (solo admin)
func getSystemInfo(c *gin.Context) {
	// TODO: Implementar lógica de información del sistema
	c.JSON(200, gin.H{
		"message": "System info - Coming soon",
	})
}

// getHealthCheck verifica el estado del sistema (solo admin)
func getHealthCheck(c *gin.Context) {
	// TODO: Implementar lógica de health check completo
	c.JSON(200, gin.H{
		"status":   "ok",
		"database": "connected",
		"message":  "System health check - Coming soon",
	})
}