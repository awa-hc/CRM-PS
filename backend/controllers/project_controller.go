package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type ProjectController struct{}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

// @Summary Obtener todos los proyectos
// @Description Obtener lista de proyectos con paginación y filtros
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Param search query string false "Buscar por nombre o descripción"
// @Param status query string false "Filtrar por estado"
// @Param client_id query int false "Filtrar por cliente"
// @Param priority query string false "Filtrar por prioridad"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /projects [get]
func (pc *ProjectController) GetProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")
	clientIDStr := c.Query("client_id")
	priority := c.Query("priority")

	offset := (page - 1) * limit

	var projects []models.Project
	var total int64

	query := config.DB.Model(&models.Project{}).Preload("Client")

	// Aplicar filtros
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if clientIDStr != "" {
		clientID, _ := strconv.ParseUint(clientIDStr, 10, 32)
		query = query.Where("client_id = ?", uint(clientID))
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Contar total
	query.Count(&total)

	// Obtener proyectos con paginación
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener proyectos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + int64(limit) - 1) / int64(limit),
	})
}

// @Summary Obtener proyecto por ID
// @Description Obtener información detallada de un proyecto
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del proyecto"
// @Success 200 {object} models.Project
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [get]
func (pc *ProjectController) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var project models.Project
	if err := config.DB.Preload("Client").First(&project, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proyecto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// @Summary Crear nuevo proyecto
// @Description Crear un nuevo proyecto
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project body models.CreateProjectRequest true "Datos del proyecto"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]string
// @Router /projects [post]
func (pc *ProjectController) CreateProject(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el cliente existe
	var client models.Client
	if err := config.DB.First(&client, req.ClientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cliente no encontrado"})
		return
	}

	// Generar código único del proyecto
	projectCode := pc.generateProjectCode()

	project := models.Project{
		Name:        req.Name,
		Description: req.Description,
		ClientID:    req.ClientID,
		Status:      "planning",
		Priority:    req.Priority,
		Type:        req.Type,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		ZipCode:     req.ZipCode,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Budget:      req.Budget,
		EstimatedCost: req.EstimatedCost,
		Progress:    0,
		Notes:       req.Notes,
		Code:        projectCode,
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear proyecto"})
		return
	}

	// Cargar el cliente para la respuesta
	config.DB.Preload("Client").First(&project, project.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Proyecto creado exitosamente",
		"project": project,
	})
}

// @Summary Actualizar proyecto
// @Description Actualizar información de un proyecto
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del proyecto"
// @Param project body models.UpdateProjectRequest true "Datos actualizados del proyecto"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Router /projects/{id} [put]
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project models.Project
	if err := config.DB.First(&project, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proyecto no encontrado"})
		return
	}

	// Actualizar campos
	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Status != nil {
		project.Status = *req.Status
	}
	if req.Priority != nil {
		project.Priority = *req.Priority
	}
	if req.Type != nil {
		project.Type = *req.Type
	}
	if req.Address != nil {
		project.Address = *req.Address
	}
	if req.City != nil {
		project.City = *req.City
	}
	if req.State != nil {
		project.State = *req.State
	}
	if req.ZipCode != nil {
		project.ZipCode = *req.ZipCode
	}
	if req.StartDate != nil {
		project.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		project.EndDate = req.EndDate
	}
	if req.Budget != nil {
		project.Budget = *req.Budget
	}
	if req.EstimatedCost != nil {
		project.EstimatedCost = *req.EstimatedCost
	}
	if req.ActualCost != nil {
		project.ActualCost = *req.ActualCost
	}
	if req.Progress != nil {
		project.Progress = *req.Progress
	}
	if req.Notes != nil {
		project.Notes = *req.Notes
	}

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar proyecto"})
		return
	}

	// Cargar el cliente para la respuesta
	config.DB.Preload("Client").First(&project, project.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Proyecto actualizado exitosamente",
		"project": project,
	})
}

// @Summary Eliminar proyecto
// @Description Eliminar un proyecto (soft delete)
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del proyecto"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [delete]
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var project models.Project
	if err := config.DB.First(&project, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proyecto no encontrado"})
		return
	}

	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar proyecto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proyecto eliminado exitosamente"})
}

// @Summary Obtener estadísticas de proyectos
// @Description Obtener estadísticas generales de proyectos
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ProjectStats
// @Router /projects/stats [get]
func (pc *ProjectController) GetProjectStats(c *gin.Context) {
	var stats models.ProjectStats

	config.DB.Model(&models.Project{}).Count(&stats.TotalProjects)
	config.DB.Model(&models.Project{}).Where("status = ?", "active").Count(&stats.ActiveProjects)
	config.DB.Model(&models.Project{}).Where("status = ?", "completed").Count(&stats.CompletedProjects)
	config.DB.Model(&models.Project{}).Where("status = ?", "planning").Count(&stats.PlanningProjects)

	// Proyectos por prioridad
	var priorityStats []struct {
		Priority string `json:"priority"`
		Count    int64  `json:"count"`
	}
	config.DB.Model(&models.Project{}).Select("priority, COUNT(*) as count").Group("priority").Scan(&priorityStats)

	// Proyectos por tipo
	var typeStats []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	config.DB.Model(&models.Project{}).Select("type, COUNT(*) as count").Group("type").Scan(&typeStats)

	// Calcular presupuesto total y costo real
	var budgetSum, actualCostSum float64
	config.DB.Model(&models.Project{}).Select("COALESCE(SUM(budget), 0)").Scan(&budgetSum)
	config.DB.Model(&models.Project{}).Select("COALESCE(SUM(actual_cost), 0)").Scan(&actualCostSum)

	stats.TotalBudget = budgetSum
	stats.TotalActualCost = actualCostSum

	// Progreso promedio
	var avgProgress float64
	config.DB.Model(&models.Project{}).Select("COALESCE(AVG(progress), 0)").Scan(&avgProgress)
	stats.AverageProgress = avgProgress

	c.JSON(http.StatusOK, gin.H{
		"stats":          stats,
		"priority_stats": priorityStats,
		"type_stats":     typeStats,
	})
}

// @Summary Obtener materiales del proyecto
// @Description Obtener lista de materiales asociados a un proyecto
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del proyecto"
// @Success 200 {object} map[string]interface{}
// @Router /projects/{id}/materials [get]
func (pc *ProjectController) GetProjectMaterials(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var projectMaterials []models.ProjectMaterial
	if err := config.DB.Preload("Material").Where("project_id = ?", uint(id)).Find(&projectMaterials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener materiales del proyecto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"materials": projectMaterials})
}

// Helper function para generar código único del proyecto
func (pc *ProjectController) generateProjectCode() string {
	now := time.Now()
	return "PRJ-" + now.Format("20060102") + "-" + strconv.FormatInt(now.Unix()%10000, 10)
}