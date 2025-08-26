package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type QuoteController struct{}

func NewQuoteController() *QuoteController {
	return &QuoteController{}
}

// @Summary Obtener todas las cotizaciones
// @Description Obtener lista de cotizaciones con paginación y filtros
// @Tags quotes
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Param status query string false "Filtrar por estado"
// @Param client_id query int false "Filtrar por cliente"
// @Success 200 {object} map[string]interface{}
// @Router /quotes [get]
func (qc *QuoteController) GetQuotes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	clientID := c.Query("client_id")

	offset := (page - 1) * limit

	query := config.DB.Model(&models.Quote{}).Preload("Client").Preload("Project").Preload("QuoteItems")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}

	var quotes []models.Quote
	var total int64

	query.Count(&total)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&quotes)

	c.JSON(http.StatusOK, gin.H{
		"quotes": quotes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// @Summary Obtener cotización por ID
// @Description Obtener detalles de una cotización específica
// @Tags quotes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la cotización"
// @Success 200 {object} models.Quote
// @Failure 404 {object} map[string]string
// @Router /quotes/{id} [get]
func (qc *QuoteController) GetQuote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var quote models.Quote
	if err := config.DB.Preload("Client").Preload("Project").Preload("QuoteItems").First(&quote, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// @Summary Crear nueva cotización
// @Description Crear una nueva cotización
// @Tags quotes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param quote body models.CreateQuoteRequest true "Datos de la cotización"
// @Success 201 {object} models.Quote
// @Failure 400 {object} map[string]string
// @Router /quotes [post]
func (qc *QuoteController) CreateQuote(c *gin.Context) {
	var req models.CreateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generar número de cotización
	quoteNumber := qc.generateQuoteNumber()

	// Calcular totales
	subtotal := 0.0
	for _, item := range req.Items {
		subtotal += item.Quantity * item.UnitPrice
	}

	taxAmount := subtotal * (req.TaxRate / 100)
	total := subtotal + taxAmount - req.Discount

	quote := models.Quote{
		QuoteNumber: quoteNumber,
		ClientID:    req.ClientID,
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "draft",
		ValidUntil:  req.ValidUntil,
		Subtotal:    subtotal,
		TaxRate:     req.TaxRate,
		TaxAmount:   taxAmount,
		Discount:    req.Discount,
		Total:       total,
		Notes:       req.Notes,
		Terms:       req.Terms,
	}

	if err := config.DB.Create(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cotización"})
		return
	}

	// Crear items de la cotización
	for _, itemReq := range req.Items {
		item := models.QuoteItem{
			QuoteID:     quote.ID,
			Description: itemReq.Description,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
			Total:       itemReq.Quantity * itemReq.UnitPrice,
		}
		config.DB.Create(&item)
	}

	// Cargar la cotización completa
	config.DB.Preload("Client").Preload("Project").Preload("QuoteItems").First(&quote, quote.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cotización creada exitosamente",
		"quote":   quote,
	})
}

// @Summary Actualizar cotización
// @Description Actualizar información de una cotización
// @Tags quotes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la cotización"
// @Param quote body models.UpdateQuoteRequest true "Datos actualizados de la cotización"
// @Success 200 {object} models.Quote
// @Failure 400 {object} map[string]string
// @Router /quotes/{id} [put]
func (qc *QuoteController) UpdateQuote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var quote models.Quote
	if err := config.DB.First(&quote, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		return
	}

	// Actualizar campos
	if req.Title != nil {
		quote.Title = *req.Title
	}
	if req.Description != nil {
		quote.Description = *req.Description
	}
	if req.Status != nil {
		quote.Status = *req.Status
	}
	if req.ValidUntil != nil {
		quote.ValidUntil = req.ValidUntil
	}
	if req.TaxRate != nil {
		quote.TaxRate = *req.TaxRate
	}
	if req.Discount != nil {
		quote.Discount = *req.Discount
	}
	if req.Notes != nil {
		quote.Notes = *req.Notes
	}
	if req.Terms != nil {
		quote.Terms = *req.Terms
	}

	// Recalcular totales si es necesario
	if req.TaxRate != nil || req.Discount != nil {
		taxAmount := quote.Subtotal * (quote.TaxRate / 100)
		total := quote.Subtotal + taxAmount - quote.Discount
		quote.TaxAmount = taxAmount
		quote.Total = total
	}

	if err := config.DB.Save(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cotización"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cotización actualizada exitosamente",
		"quote":   quote,
	})
}

// @Summary Eliminar cotización
// @Description Eliminar una cotización
// @Tags quotes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la cotización"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /quotes/{id} [delete]
func (qc *QuoteController) DeleteQuote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var quote models.Quote
	if err := config.DB.First(&quote, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		return
	}

	// Verificar si la cotización puede ser eliminada
	if quote.Status == "accepted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar una cotización aceptada"})
		return
	}

	// Eliminar items de la cotización
	config.DB.Where("quote_id = ?", quote.ID).Delete(&models.QuoteItem{})

	if err := config.DB.Delete(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cotización"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cotización eliminada exitosamente"})
}

// @Summary Obtener estadísticas de cotizaciones
// @Description Obtener estadísticas generales de cotizaciones
// @Tags quotes
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /quotes/stats [get]
func (qc *QuoteController) GetQuoteStats(c *gin.Context) {
	var stats models.QuoteStats

	config.DB.Model(&models.Quote{}).Count(&stats.TotalQuotes)
	config.DB.Model(&models.Quote{}).Where("status = ?", "draft").Count(&stats.DraftQuotes)
	config.DB.Model(&models.Quote{}).Where("status = ?", "sent").Count(&stats.SentQuotes)
	config.DB.Model(&models.Quote{}).Where("status = ?", "accepted").Count(&stats.AcceptedQuotes)
	config.DB.Model(&models.Quote{}).Where("status = ?", "rejected").Count(&stats.RejectedQuotes)

	// Valor total de cotizaciones
	var totalValue float64
	config.DB.Model(&models.Quote{}).Select("COALESCE(SUM(total), 0)").Scan(&totalValue)
	stats.TotalValue = totalValue

	// Valor de cotizaciones aceptadas
	var acceptedValue float64
	config.DB.Model(&models.Quote{}).Where("status = ?", "accepted").Select("COALESCE(SUM(total), 0)").Scan(&acceptedValue)
	stats.AcceptedValue = acceptedValue

	// Cotizaciones del mes actual
	var thisMonthQuotes int64
	config.DB.Model(&models.Quote{}).Where("created_at >= DATE_TRUNC('month', CURRENT_DATE)").Count(&thisMonthQuotes)
	stats.ThisMonthQuotes = thisMonthQuotes

	c.JSON(http.StatusOK, stats)
}

// Función auxiliar para generar número de cotización
func (qc *QuoteController) generateQuoteNumber() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	// Contar cotizaciones del mes actual
	var count int64
	config.DB.Model(&models.Quote{}).Where("created_at >= DATE_TRUNC('month', CURRENT_DATE)").Count(&count)

	return fmt.Sprintf("COT-%d%02d-%04d", year, month, count+1)
}

// @Summary Cambiar estado de cotización
// @Description Cambiar el estado de una cotización (draft, sent, accepted, rejected)
// @Tags quotes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la cotización"
// @Param status body map[string]string true "Nuevo estado"
// @Success 200 {object} models.Quote
// @Failure 400 {object} map[string]string
// @Router /quotes/{id}/status [patch]
func (qc *QuoteController) ChangeQuoteStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar estado
	validStatuses := []string{"draft", "sent", "accepted", "rejected"}
	valid := false
	for _, status := range validStatuses {
		if req.Status == status {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido"})
		return
	}

	var quote models.Quote
	if err := config.DB.First(&quote, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		return
	}

	quote.Status = req.Status

	if err := config.DB.Save(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar estado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estado actualizado exitosamente",
		"quote":   quote,
	})
}