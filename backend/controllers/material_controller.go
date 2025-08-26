package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type MaterialController struct{}

func NewMaterialController() *MaterialController {
	return &MaterialController{}
}

// @Summary Obtener todos los materiales
// @Description Obtener lista de materiales con paginación y filtros
// @Tags materials
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Param category query string false "Filtrar por categoría"
// @Param search query string false "Buscar por nombre o descripción"
// @Success 200 {object} map[string]interface{}
// @Router /materials [get]
func (mc *MaterialController) GetMaterials(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	search := c.Query("search")

	offset := (page - 1) * limit

	query := config.DB.Model(&models.Material{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var materials []models.Material
	var total int64

	query.Count(&total)
	query.Offset(offset).Limit(limit).Order("name ASC").Find(&materials)

	c.JSON(http.StatusOK, gin.H{
		"materials": materials,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// @Summary Obtener material por ID
// @Description Obtener detalles de un material específico
// @Tags materials
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del material"
// @Success 200 {object} models.Material
// @Failure 404 {object} map[string]string
// @Router /materials/{id} [get]
func (mc *MaterialController) GetMaterial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var material models.Material
	if err := config.DB.First(&material, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material no encontrado"})
		return
	}

	c.JSON(http.StatusOK, material)
}

// @Summary Crear nuevo material
// @Description Crear un nuevo material en el inventario
// @Tags materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param material body models.CreateMaterialRequest true "Datos del material"
// @Success 201 {object} models.Material
// @Failure 400 {object} map[string]string
// @Router /materials [post]
func (mc *MaterialController) CreateMaterial(c *gin.Context) {
	var req models.CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si ya existe un material con el mismo SKU
	if req.SKU != "" {
		var existingMaterial models.Material
		if err := config.DB.Where("sku = ?", req.SKU).First(&existingMaterial).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ya existe un material con este SKU"})
			return
		}
	}

	material := models.Material{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Unit:        req.Unit,
		UnitPrice:   req.UnitPrice,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		Supplier:    req.Supplier,
		SKU:         req.SKU,
		IsActive:    true,
	}

	if err := config.DB.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear material"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Material creado exitosamente",
		"material": material,
	})
}

// @Summary Actualizar material
// @Description Actualizar información de un material
// @Tags materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del material"
// @Param material body models.UpdateMaterialRequest true "Datos actualizados del material"
// @Success 200 {object} models.Material
// @Failure 400 {object} map[string]string
// @Router /materials/{id} [put]
func (mc *MaterialController) UpdateMaterial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var material models.Material
	if err := config.DB.First(&material, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material no encontrado"})
		return
	}

	// Verificar si el SKU ya existe en otro material
	if req.SKU != nil && *req.SKU != material.SKU {
		var existingMaterial models.Material
		if err := config.DB.Where("sku = ? AND id != ?", *req.SKU, material.ID).First(&existingMaterial).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ya existe otro material con este SKU"})
			return
		}
	}

	// Actualizar campos
	if req.Name != nil {
		material.Name = *req.Name
	}
	if req.Description != nil {
		material.Description = *req.Description
	}
	if req.Category != nil {
		material.Category = *req.Category
	}
	if req.Unit != nil {
		material.Unit = *req.Unit
	}
	if req.UnitPrice != nil {
		material.UnitPrice = *req.UnitPrice
	}
	if req.Stock != nil {
		material.Stock = *req.Stock
	}
	if req.MinStock != nil {
		material.MinStock = *req.MinStock
	}
	if req.Supplier != nil {
		material.Supplier = *req.Supplier
	}
	if req.SKU != nil {
		material.SKU = *req.SKU
	}
	if req.IsActive != nil {
		material.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar material"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Material actualizado exitosamente",
		"material": material,
	})
}

// @Summary Eliminar material
// @Description Eliminar un material del inventario
// @Tags materials
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del material"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /materials/{id} [delete]
func (mc *MaterialController) DeleteMaterial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var material models.Material
	if err := config.DB.First(&material, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material no encontrado"})
		return
	}

	// Verificar si el material está siendo usado en proyectos
	var projectMaterialCount int64
	config.DB.Model(&models.ProjectMaterial{}).Where("material_id = ?", material.ID).Count(&projectMaterialCount)

	if projectMaterialCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede eliminar el material porque está siendo usado en proyectos"})
		return
	}

	if err := config.DB.Delete(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar material"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material eliminado exitosamente"})
}

// @Summary Actualizar stock de material
// @Description Actualizar el stock de un material (entrada o salida)
// @Tags materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del material"
// @Param stock body map[string]interface{} true "Datos de actualización de stock"
// @Success 200 {object} models.Material
// @Failure 400 {object} map[string]string
// @Router /materials/{id}/stock [patch]
func (mc *MaterialController) UpdateMaterialStock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		Quantity float64 `json:"quantity" binding:"required"`
		Type     string  `json:"type" binding:"required"` // "in" o "out"
		Reason   string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type != "in" && req.Type != "out" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo debe ser 'in' o 'out'"})
		return
	}

	var material models.Material
	if err := config.DB.First(&material, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material no encontrado"})
		return
	}

	// Actualizar stock
	if req.Type == "in" {
		material.Stock += req.Quantity
	} else {
		if material.Stock < req.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente"})
			return
		}
		material.Stock -= req.Quantity
	}

	if err := config.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Stock actualizado exitosamente",
		"material": material,
	})
}

// @Summary Obtener materiales con stock bajo
// @Description Obtener lista de materiales con stock por debajo del mínimo
// @Tags materials
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /materials/low-stock [get]
func (mc *MaterialController) GetLowStockMaterials(c *gin.Context) {
	var materials []models.Material
	config.DB.Where("stock <= min_stock AND is_active = true").Find(&materials)

	c.JSON(http.StatusOK, gin.H{
		"materials": materials,
		"count":     len(materials),
	})
}

// @Summary Obtener categorías de materiales
// @Description Obtener lista de categorías únicas de materiales
// @Tags materials
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /materials/categories [get]
func (mc *MaterialController) GetMaterialCategories(c *gin.Context) {
	var categories []string
	config.DB.Model(&models.Material{}).Distinct("category").Where("category != ''").Pluck("category", &categories)

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// @Summary Obtener estadísticas de materiales
// @Description Obtener estadísticas generales del inventario
// @Tags materials
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /materials/stats [get]
func (mc *MaterialController) GetMaterialStats(c *gin.Context) {
	var totalMaterials, activeMaterials, lowStockMaterials int64
	var totalValue float64

	config.DB.Model(&models.Material{}).Count(&totalMaterials)
	config.DB.Model(&models.Material{}).Where("is_active = true").Count(&activeMaterials)
	config.DB.Model(&models.Material{}).Where("stock <= min_stock AND is_active = true").Count(&lowStockMaterials)
	config.DB.Model(&models.Material{}).Select("COALESCE(SUM(stock * unit_price), 0)").Scan(&totalValue)

	// Materiales por categoría
	var categoriesStats []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	config.DB.Model(&models.Material{}).Select("category, COUNT(*) as count").Where("is_active = true").Group("category").Scan(&categoriesStats)

	c.JSON(http.StatusOK, gin.H{
		"total_materials":     totalMaterials,
		"active_materials":    activeMaterials,
		"low_stock_materials": lowStockMaterials,
		"total_value":         totalValue,
		"categories_stats":    categoriesStats,
	})
}