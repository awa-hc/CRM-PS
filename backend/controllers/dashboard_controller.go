package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// @Summary Obtener estadísticas del dashboard
// @Description Obtener estadísticas generales para el dashboard principal
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /dashboard/stats [get]
func (dc *DashboardController) GetDashboardStats(c *gin.Context) {
	stats := make(map[string]interface{})

	// Estadísticas de clientes
	var totalClients, activeClients int64
	config.DB.Model(&models.Client{}).Count(&totalClients)
	config.DB.Model(&models.Client{}).Where("is_active = true").Count(&activeClients)

	stats["clients"] = map[string]interface{}{
		"total":  totalClients,
		"active": activeClients,
	}

	// Estadísticas de proyectos
	var totalProjects, activeProjects, completedProjects int64
	config.DB.Model(&models.Project{}).Count(&totalProjects)
	config.DB.Model(&models.Project{}).Where("status IN ?", []string{"planning", "in_progress", "on_hold"}).Count(&activeProjects)
	config.DB.Model(&models.Project{}).Where("status = ?", "completed").Count(&completedProjects)

	stats["projects"] = map[string]interface{}{
		"total":     totalProjects,
		"active":    activeProjects,
		"completed": completedProjects,
	}

	// Estadísticas de cotizaciones
	var totalQuotes, pendingQuotes, acceptedQuotes int64
	var quotesValue, acceptedValue float64
	config.DB.Model(&models.Quote{}).Count(&totalQuotes)
	config.DB.Model(&models.Quote{}).Where("status IN ?", []string{"draft", "sent"}).Count(&pendingQuotes)
	config.DB.Model(&models.Quote{}).Where("status = ?", "accepted").Count(&acceptedQuotes)
	config.DB.Model(&models.Quote{}).Select("COALESCE(SUM(total), 0)").Scan(&quotesValue)
	config.DB.Model(&models.Quote{}).Where("status = ?", "accepted").Select("COALESCE(SUM(total), 0)").Scan(&acceptedValue)

	stats["quotes"] = map[string]interface{}{
		"total":          totalQuotes,
		"pending":        pendingQuotes,
		"accepted":       acceptedQuotes,
		"total_value":    quotesValue,
		"accepted_value": acceptedValue,
	}

	// Estadísticas de materiales
	var totalMaterials, lowStockMaterials int64
	var inventoryValue float64
	config.DB.Model(&models.Material{}).Where("is_active = true").Count(&totalMaterials)
	config.DB.Model(&models.Material{}).Where("stock <= min_stock AND is_active = true").Count(&lowStockMaterials)
	config.DB.Model(&models.Material{}).Select("COALESCE(SUM(stock * price), 0)").Scan(&inventoryValue)

	stats["materials"] = map[string]interface{}{
		"total":           totalMaterials,
		"low_stock":       lowStockMaterials,
		"inventory_value": inventoryValue,
	}

	// Estadísticas mensuales
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var monthlyClients, monthlyProjects, monthlyQuotes int64
	config.DB.Model(&models.Client{}).Where("created_at >= ?", startOfMonth).Count(&monthlyClients)
	config.DB.Model(&models.Project{}).Where("created_at >= ?", startOfMonth).Count(&monthlyProjects)
	config.DB.Model(&models.Quote{}).Where("created_at >= ?", startOfMonth).Count(&monthlyQuotes)

	stats["monthly"] = map[string]interface{}{
		"clients":  monthlyClients,
		"projects": monthlyProjects,
		"quotes":   monthlyQuotes,
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Obtener actividad reciente
// @Description Obtener lista de actividades recientes en el sistema
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Número de elementos" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /dashboard/recent-activity [get]
func (dc *DashboardController) GetRecentActivity(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	activities := []map[string]interface{}{}

	// Clientes recientes
	var recentClients []models.Client
	config.DB.Order("created_at DESC").Limit(limit/4 + 1).Find(&recentClients)
	for _, client := range recentClients {
		activities = append(activities, map[string]interface{}{
			"type":        "client",
			"action":      "created",
			"description": "Nuevo cliente: " + client.Name,
			"created_at":  client.CreatedAt,
			"id":          client.ID,
		})
	}

	// Proyectos recientes
	var recentProjects []models.Project
	config.DB.Preload("Client").Order("created_at DESC").Limit(limit/4 + 1).Find(&recentProjects)
	for _, project := range recentProjects {
		clientName := "Cliente desconocido"
		if project.Client.ID != 0 {
			clientName = project.Client.Name
		}
		activities = append(activities, map[string]interface{}{
			"type":        "project",
			"action":      "created",
			"description": "Nuevo proyecto: " + project.Name + " para " + clientName,
			"created_at":  project.CreatedAt,
			"id":          project.ID,
		})
	}

	// Cotizaciones recientes
	var recentQuotes []models.Quote
	config.DB.Preload("Client").Order("created_at DESC").Limit(limit/4 + 1).Find(&recentQuotes)
	for _, quote := range recentQuotes {
		clientName := "Cliente desconocido"
		if quote.Client.ID != 0 {
			clientName = quote.Client.Name
		}
		activities = append(activities, map[string]interface{}{
			"type":        "quote",
			"action":      "created",
			"description": "Nueva cotización: " + quote.Title + " para " + clientName,
			"created_at":  quote.CreatedAt,
			"id":          quote.ID,
		})
	}

	// Materiales con stock bajo
	var lowStockMaterials []models.Material
	config.DB.Where("stock <= min_stock AND is_active = true").Limit(limit/4 + 1).Find(&lowStockMaterials)
	for _, material := range lowStockMaterials {
		activities = append(activities, map[string]interface{}{
			"type":        "material",
			"action":      "low_stock",
			"description": "Stock bajo: " + material.Name + " (" + strconv.FormatFloat(material.Stock, 'f', 2, 64) + " " + material.Unit + ")",
			"created_at":  material.UpdatedAt,
			"id":          material.ID,
		})
	}

	// Ordenar por fecha y limitar
	sort.Slice(activities, func(i, j int) bool {
		timeI := activities[i]["created_at"].(time.Time)
		timeJ := activities[j]["created_at"].(time.Time)
		return timeI.After(timeJ)
	})

	if len(activities) > limit {
		activities = activities[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"count":      len(activities),
	})
}

// @Summary Obtener gráfico de proyectos por estado
// @Description Obtener datos para gráfico de proyectos agrupados por estado
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /dashboard/projects-by-status [get]
func (dc *DashboardController) GetProjectsByStatus(c *gin.Context) {
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	config.DB.Model(&models.Project{}).Select("status, COUNT(*) as count").Group("status").Scan(&statusStats)

	c.JSON(http.StatusOK, gin.H{
		"data": statusStats,
	})
}

// @Summary Obtener gráfico de ingresos mensuales
// @Description Obtener datos de ingresos de los últimos 12 meses
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /dashboard/monthly-revenue [get]
func (dc *DashboardController) GetMonthlyRevenue(c *gin.Context) {
	now := time.Now()
	monthlyData := []map[string]interface{}{}

	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		var revenue float64
		config.DB.Model(&models.Quote{}).Where("status = ? AND created_at BETWEEN ? AND ?", "accepted", startOfMonth, endOfMonth).Select("COALESCE(SUM(total), 0)").Scan(&revenue)

		monthlyData = append(monthlyData, map[string]interface{}{
			"month":   month.Format("2006-01"),
			"revenue": revenue,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": monthlyData,
	})
}

// @Summary Obtener proyectos próximos a vencer
// @Description Obtener lista de proyectos que están próximos a su fecha de finalización
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /dashboard/upcoming-deadlines [get]
func (dc *DashboardController) GetUpcomingDeadlines(c *gin.Context) {
	now := time.Now()
	next30Days := now.AddDate(0, 0, 30)

	var projects []models.Project
	config.DB.Preload("Client").Where("end_date BETWEEN ? AND ? AND status IN ?", now, next30Days, []string{"planning", "in_progress"}).Order("end_date ASC").Find(&projects)

	projectData := []map[string]interface{}{}
	for _, project := range projects {
		clientName := "Cliente desconocido"
		if project.Client.ID != 0 {
			clientName = project.Client.Name
		}

		daysLeft := int(project.EndDate.Sub(now).Hours() / 24)
		if daysLeft < 0 {
			daysLeft = 0
		}

		projectData = append(projectData, map[string]interface{}{
			"id":          project.ID,
			"name":        project.Name,
			"client":      clientName,
			"end_date":    project.EndDate,
			"days_left":   daysLeft,
			"status":      project.Status,
			"progress":    project.Progress,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projectData,
		"count":    len(projectData),
	})
}

// @Summary Obtener resumen financiero
// @Description Obtener resumen financiero del mes actual
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /dashboard/financial-summary [get]
func (dc *DashboardController) GetFinancialSummary(c *gin.Context) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfMonth.AddDate(0, 0, -1)

	// Ingresos del mes actual
	var currentMonthRevenue float64
	config.DB.Model(&models.Quote{}).Where("status = ? AND created_at >= ?", "accepted", startOfMonth).Select("COALESCE(SUM(total), 0)").Scan(&currentMonthRevenue)

	// Ingresos del mes pasado
	var lastMonthRevenue float64
	config.DB.Model(&models.Quote{}).Where("status = ? AND created_at BETWEEN ? AND ?", "accepted", startOfLastMonth, endOfLastMonth).Select("COALESCE(SUM(total), 0)").Scan(&lastMonthRevenue)

	// Calcular crecimiento
	var growth float64
	if lastMonthRevenue > 0 {
		growth = ((currentMonthRevenue - lastMonthRevenue) / lastMonthRevenue) * 100
	}

	// Cotizaciones pendientes
	var pendingQuotesValue float64
	config.DB.Model(&models.Quote{}).Where("status IN ?", []string{"draft", "sent"}).Select("COALESCE(SUM(total), 0)").Scan(&pendingQuotesValue)

	// Valor del inventario
	var inventoryValue float64
	config.DB.Model(&models.Material{}).Select("COALESCE(SUM(stock * price), 0)").Scan(&inventoryValue)

	c.JSON(http.StatusOK, gin.H{
		"current_month_revenue": currentMonthRevenue,
		"last_month_revenue":    lastMonthRevenue,
		"growth_percentage":     growth,
		"pending_quotes_value":  pendingQuotesValue,
		"inventory_value":       inventoryValue,
	})
}