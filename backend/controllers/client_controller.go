package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type ClientController struct{}

func NewClientController() *ClientController {
	return &ClientController{}
}

// @Summary Obtener todos los clientes
// @Description Obtener lista de clientes con paginación y filtros
// @Tags clients
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Param search query string false "Buscar por nombre o email"
// @Param active query bool false "Filtrar por estado activo"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /clients [get]
func (cc *ClientController) GetClients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	activeStr := c.Query("active")

	offset := (page - 1) * limit

	var clients []models.Client
	var total int64

	query := config.DB.Model(&models.Client{})

	// Aplicar filtros
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR company ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if activeStr != "" {
		active, _ := strconv.ParseBool(activeStr)
		query = query.Where("is_active = ?", active)
	}

	// Contar total
	query.Count(&total)

	// Obtener clientes con paginación
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener clientes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"clients": clients,
		"total":   total,
		"page":    page,
		"limit":   limit,
		"pages":   (total + int64(limit) - 1) / int64(limit),
	})
}

// @Summary Obtener cliente por ID
// @Description Obtener información detallada de un cliente
// @Tags clients
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del cliente"
// @Success 200 {object} models.Client
// @Failure 404 {object} map[string]string
// @Router /clients/{id} [get]
func (cc *ClientController) GetClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var client models.Client
	if err := config.DB.First(&client, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// @Summary Crear nuevo cliente
// @Description Crear un nuevo cliente
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param client body models.CreateClientRequest true "Datos del cliente"
// @Success 201 {object} models.Client
// @Failure 400 {object} map[string]string
// @Router /clients [post]
func (cc *ClientController) CreateClient(c *gin.Context) {
	var req models.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si ya existe un cliente con el mismo email
	var existingClient models.Client
	if err := config.DB.Where("email = ?", req.Email).First(&existingClient).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ya existe un cliente con este email"})
		return
	}

	client := models.Client{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		ZipCode:     req.ZipCode,
		Company:     req.Company,
		TaxID:       req.TaxID,
		ContactType: req.ContactType,
		Notes:       req.Notes,
		IsActive:    true,
	}

	if err := config.DB.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cliente creado exitosamente",
		"client":  client,
	})
}

// @Summary Actualizar cliente
// @Description Actualizar información de un cliente
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del cliente"
// @Param client body models.UpdateClientRequest true "Datos actualizados del cliente"
// @Success 200 {object} models.Client
// @Failure 400 {object} map[string]string
// @Router /clients/{id} [put]
func (cc *ClientController) UpdateClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var client models.Client
	if err := config.DB.First(&client, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	// Verificar si el email ya existe en otro cliente
	if req.Email != nil && *req.Email != client.Email {
		var existingClient models.Client
		if err := config.DB.Where("email = ? AND id != ?", *req.Email, client.ID).First(&existingClient).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ya existe otro cliente con este email"})
			return
		}
	}

	// Actualizar campos
	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.Email != nil {
		client.Email = *req.Email
	}
	if req.Phone != nil {
		client.Phone = *req.Phone
	}
	if req.Address != nil {
		client.Address = *req.Address
	}
	if req.City != nil {
		client.City = *req.City
	}
	if req.State != nil {
		client.State = *req.State
	}
	if req.ZipCode != nil {
		client.ZipCode = *req.ZipCode
	}
	if req.Company != nil {
		client.Company = *req.Company
	}
	if req.TaxID != nil {
		client.TaxID = *req.TaxID
	}
	if req.ContactType != nil {
		client.ContactType = *req.ContactType
	}
	if req.Notes != nil {
		client.Notes = *req.Notes
	}
	if req.IsActive != nil {
		client.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cliente actualizado exitosamente",
		"client":  client,
	})
}

// @Summary Eliminar cliente
// @Description Eliminar un cliente (soft delete)
// @Tags clients
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del cliente"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clients/{id} [delete]
func (cc *ClientController) DeleteClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var client models.Client
	if err := config.DB.First(&client, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	// Verificar si el cliente tiene proyectos activos
	var projectCount int64
	config.DB.Model(&models.Project{}).Where("client_id = ? AND status NOT IN ?", client.ID, []string{"completed", "cancelled"}).Count(&projectCount)

	if projectCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar el cliente porque tiene proyectos activos"})
		return
	}

	if err := config.DB.Delete(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado exitosamente"})
}

// @Summary Obtener estadísticas de clientes
// @Description Obtener estadísticas generales de clientes
// @Tags clients
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /clients/stats [get]
func (cc *ClientController) GetClientStats(c *gin.Context) {
	var totalClients, activeClients, inactiveClients int64

	config.DB.Model(&models.Client{}).Count(&totalClients)
	config.DB.Model(&models.Client{}).Where("is_active = ?", true).Count(&activeClients)
	config.DB.Model(&models.Client{}).Where("is_active = ?", false).Count(&inactiveClients)

	// Clientes por tipo de contacto
	var contactTypes []struct {
		ContactType string `json:"contact_type"`
		Count       int64  `json:"count"`
	}
	config.DB.Model(&models.Client{}).Select("contact_type, COUNT(*) as count").Group("contact_type").Scan(&contactTypes)

	// Clientes registrados en los últimos 30 días
	var recentClients int64
	config.DB.Model(&models.Client{}).Where("created_at >= NOW() - INTERVAL '30 days'").Count(&recentClients)

	c.JSON(http.StatusOK, gin.H{
		"total_clients":    totalClients,
		"active_clients":   activeClients,
		"inactive_clients": inactiveClients,
		"contact_types":    contactTypes,
		"recent_clients":   recentClients,
	})
}