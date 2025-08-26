package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type ReportController struct{}

func NewReportController() *ReportController {
	return &ReportController{}
}

// @Summary Reporte de clientes
// @Description Generar reporte detallado de clientes
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Param status query string false "Filtrar por estado (active/inactive)"
// @Success 200 {object} map[string]interface{}
// @Router /reports/clients [get]
func (rc *ReportController) GetClientsReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	query := config.DB.Model(&models.Client{})

	// Aplicar filtros de fecha
	if startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", parsedDate)
		}
	}
	if endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("created_at <= ?", parsedDate.AddDate(0, 0, 1))
		}
	}

	// Aplicar filtro de estado
	if status == "active" {
		query = query.Where("is_active = true")
	} else if status == "inactive" {
		query = query.Where("is_active = false")
	}

	var clients []models.Client
	query.Preload("Projects").Preload("Quotes").Find(&clients)

	// Calcular estadísticas
	totalClients := len(clients)
	activeClients := 0
	totalProjects := 0
	totalQuotes := 0
	var totalQuoteValue float64

	clientData := []map[string]interface{}{}
	for _, client := range clients {
		if client.IsActive {
			activeClients++
		}
		totalProjects += len(client.Projects)
		totalQuotes += len(client.Quotes)

		var clientQuoteValue float64
		for _, quote := range client.Quotes {
			clientQuoteValue += quote.Total
			totalQuoteValue += quote.Total
		}

		clientData = append(clientData, map[string]interface{}{
			"id":           client.ID,
			"name":         client.Name,
			"email":        client.Email,
			"company":      client.Company,
			"contact_type": client.ContactType,
			"is_active":    client.IsActive,
			"created_at":   client.CreatedAt,
			"projects":     len(client.Projects),
			"quotes":       len(client.Quotes),
			"quote_value":  clientQuoteValue,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": map[string]interface{}{
			"total_clients":    totalClients,
			"active_clients":   activeClients,
			"inactive_clients": totalClients - activeClients,
			"total_projects":   totalProjects,
			"total_quotes":     totalQuotes,
			"total_value":      totalQuoteValue,
		},
		"clients": clientData,
	})
}

// @Summary Reporte de proyectos
// @Description Generar reporte detallado de proyectos
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Param status query string false "Filtrar por estado"
// @Param client_id query int false "Filtrar por cliente"
// @Success 200 {object} map[string]interface{}
// @Router /reports/projects [get]
func (rc *ReportController) GetProjectsReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	clientID := c.Query("client_id")

	query := config.DB.Model(&models.Project{})

	// Aplicar filtros
	if startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", parsedDate)
		}
	}
	if endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("created_at <= ?", parsedDate.AddDate(0, 0, 1))
		}
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}

	var projects []models.Project
	query.Preload("Client").Preload("Quotes").Preload("ProjectMaterials").Find(&projects)

	// Calcular estadísticas
	statusCount := make(map[string]int)
	var totalBudget, totalCost float64
	totalProjects := len(projects)

	projectData := []map[string]interface{}{}
	for _, project := range projects {
		statusCount[project.Status]++
		totalBudget += project.Budget
		totalCost += project.ActualCost

		clientName := "Cliente desconocido"
		if project.Client.ID != 0 {
			clientName = project.Client.Name
		}

		// Calcular costo de materiales
		var materialsCost float64
		for _, pm := range project.ProjectMaterials {
			materialsCost += pm.QuantityPlanned * pm.UnitPrice
		}

		projectData = append(projectData, map[string]interface{}{
			"id":             project.ID,
			"code":           project.Code,
			"name":           project.Name,
			"client":         clientName,
			"status":         project.Status,
			"priority":       project.Priority,
			"type":           project.Type,
			"budget":         project.Budget,
			"cost":           project.ActualCost,
			"materials_cost": materialsCost,
			"progress":       project.Progress,
			"start_date":     project.StartDate,
			"end_date":       project.EndDate,
			"created_at":     project.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": map[string]interface{}{
			"total_projects": totalProjects,
			"status_count":   statusCount,
			"total_budget":   totalBudget,
			"total_cost":     totalCost,
			"profit_margin":  totalBudget - totalCost,
		},
		"projects": projectData,
	})
}

// @Summary Reporte de cotizaciones
// @Description Generar reporte detallado de cotizaciones
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Param status query string false "Filtrar por estado"
// @Param client_id query int false "Filtrar por cliente"
// @Success 200 {object} map[string]interface{}
// @Router /reports/quotes [get]
func (rc *ReportController) GetQuotesReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	clientID := c.Query("client_id")

	query := config.DB.Model(&models.Quote{})

	// Aplicar filtros
	if startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", parsedDate)
		}
	}
	if endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("created_at <= ?", parsedDate.AddDate(0, 0, 1))
		}
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}

	var quotes []models.Quote
	query.Preload("Client").Preload("Project").Preload("Items").Find(&quotes)

	// Calcular estadísticas
	statusCount := make(map[string]int)
	statusValue := make(map[string]float64)
	var totalValue float64
	totalQuotes := len(quotes)

	quoteData := []map[string]interface{}{}
	for _, quote := range quotes {
		statusCount[quote.Status]++
		statusValue[quote.Status] += quote.Total
		totalValue += quote.Total

		clientName := "Cliente desconocido"
		if quote.Client.ID != 0 {
			clientName = quote.Client.Name
		}

		projectName := ""
		if quote.Project.ID != 0 {
			projectName = quote.Project.Name
		}

		quoteData = append(quoteData, map[string]interface{}{
			"id":           quote.ID,
			"quote_number": quote.QuoteNumber,
			"title":        quote.Title,
			"client":       clientName,
			"project":      projectName,
			"status":       quote.Status,
			"subtotal":     quote.Subtotal,
			"tax_rate":     quote.TaxRate,
			"tax_amount":   quote.TaxAmount,
			"discount":     quote.Discount,
			"total":        quote.Total,
			"valid_until":  quote.ValidUntil,
			"created_at":   quote.CreatedAt,
			"items_count":  len(quote.Items),
		})
	}

	// Calcular tasa de conversión
	conversionRate := 0.0
	if totalQuotes > 0 {
		acceptedCount := statusCount["accepted"]
		conversionRate = (float64(acceptedCount) / float64(totalQuotes)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": map[string]interface{}{
			"total_quotes":    totalQuotes,
			"status_count":    statusCount,
			"status_value":    statusValue,
			"total_value":     totalValue,
			"conversion_rate": conversionRate,
		},
		"quotes": quoteData,
	})
}

// @Summary Reporte de materiales
// @Description Generar reporte de inventario de materiales
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param category query string false "Filtrar por categoría"
// @Param low_stock query bool false "Solo materiales con stock bajo"
// @Success 200 {object} map[string]interface{}
// @Router /reports/materials [get]
func (rc *ReportController) GetMaterialsReport(c *gin.Context) {
	category := c.Query("category")
	lowStock := c.Query("low_stock") == "true"

	query := config.DB.Model(&models.Material{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if lowStock {
		query = query.Where("stock <= min_stock")
	}

	var materials []models.Material
	query.Find(&materials)

	// Calcular estadísticas
	categoryStats := make(map[string]map[string]interface{})
	var totalValue, totalStock float64
	lowStockCount := 0

	materialData := []map[string]interface{}{}
	for _, material := range materials {
		value := material.Stock * material.UnitPrice
		totalValue += value
		totalStock += material.Stock

		if material.Stock <= material.MinStock {
			lowStockCount++
		}

		// Estadísticas por categoría
		if _, exists := categoryStats[material.Category]; !exists {
			categoryStats[material.Category] = map[string]interface{}{
				"count": 0,
				"value": 0.0,
				"stock": 0.0,
			}
		}
		categoryStats[material.Category]["count"] = categoryStats[material.Category]["count"].(int) + 1
		categoryStats[material.Category]["value"] = categoryStats[material.Category]["value"].(float64) + value
		categoryStats[material.Category]["stock"] = categoryStats[material.Category]["stock"].(float64) + material.Stock

		materialData = append(materialData, map[string]interface{}{
			"id":          material.ID,
			"name":        material.Name,
			"category":    material.Category,
			"unit":        material.Unit,
			"price":       material.UnitPrice,
			"stock":       material.Stock,
			"min_stock":   material.MinStock,
			"value":       value,
			"supplier":    material.Supplier,
			"sku":         material.SKU,
			"is_active":   material.IsActive,
			"low_stock":   material.Stock <= material.MinStock,
			"created_at":  material.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": map[string]interface{}{
			"total_materials": len(materials),
			"total_value":     totalValue,
			"total_stock":     totalStock,
			"low_stock_count": lowStockCount,
			"category_stats":  categoryStats,
		},
		"materials": materialData,
	})
}

// @Summary Reporte financiero
// @Description Generar reporte financiero consolidado
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /reports/financial [get]
func (rc *ReportController) GetFinancialReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Fechas por defecto (último mes)
	now := time.Now()
	defaultStart := now.AddDate(0, -1, 0)
	defaultEnd := now

	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			start = defaultStart
		}
	} else {
		start = defaultStart
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			end = defaultEnd
		}
	} else {
		end = defaultEnd
	}

	// Ingresos (cotizaciones aceptadas)
	var revenue float64
	config.DB.Model(&models.Quote{}).Where("status = ? AND created_at BETWEEN ? AND ?", "accepted", start, end).Select("COALESCE(SUM(total), 0)").Scan(&revenue)

	// Costos de proyectos
	var projectCosts float64
	config.DB.Model(&models.Project{}).Where("created_at BETWEEN ? AND ?", start, end).Select("COALESCE(SUM(cost), 0)").Scan(&projectCosts)

	// Valor del inventario
	var inventoryValue float64
	config.DB.Model(&models.Material{}).Select("COALESCE(SUM(stock * price), 0)").Scan(&inventoryValue)

	// Cotizaciones pendientes
	var pendingQuotes float64
	config.DB.Model(&models.Quote{}).Where("status IN ? AND created_at BETWEEN ? AND ?", []string{"draft", "sent"}, start, end).Select("COALESCE(SUM(total), 0)").Scan(&pendingQuotes)

	// Proyectos activos
	var activeProjectsValue float64
	config.DB.Model(&models.Project{}).Where("status IN ? AND created_at BETWEEN ? AND ?", []string{"planning", "in_progress"}, start, end).Select("COALESCE(SUM(budget), 0)").Scan(&activeProjectsValue)

	// Calcular métricas
	profit := revenue - projectCosts
	profitMargin := 0.0
	if revenue > 0 {
		profitMargin = (profit / revenue) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"period": map[string]interface{}{
			"start_date": start.Format("2006-01-02"),
			"end_date":   end.Format("2006-01-02"),
		},
		"revenue": map[string]interface{}{
			"total":        revenue,
			"pending":      pendingQuotes,
			"active_value": activeProjectsValue,
		},
		"costs": map[string]interface{}{
			"projects": projectCosts,
		},
		"profit": map[string]interface{}{
			"amount": profit,
			"margin": profitMargin,
		},
		"assets": map[string]interface{}{
			"inventory_value": inventoryValue,
		},
	})
}