package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"raborimet-crm/backend/models"
	"raborimet-crm/backend/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// @Summary Registrar nuevo usuario
// @Description Crear una nueva cuenta de usuario
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Datos del usuario"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado exitosamente",
		"user":    user,
	})
}

// @Summary Iniciar sesión
// @Description Autenticar usuario y obtener token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Credenciales de acceso"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := ac.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// @Summary Obtener perfil del usuario
// @Description Obtener información del usuario autenticado
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /auth/profile [get]
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	user, err := ac.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Actualizar perfil del usuario
// @Description Actualizar información del usuario autenticado
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updates body map[string]interface{} true "Campos a actualizar"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /auth/profile [put]
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.authService.UpdateUserProfile(userID.(uint), updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado exitosamente",
		"user":    user,
	})
}

// @Summary Cambiar contraseña
// @Description Cambiar la contraseña del usuario autenticado
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "Contraseñas"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/change-password [post]
func (ac *AuthController) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.authService.ChangePassword(userID.(uint), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña cambiada exitosamente"})
}

// @Summary Verificar token
// @Description Verificar si el token JWT es válido
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/verify [get]
func (ac *AuthController) VerifyToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	// Obtener el usuario completo de la base de datos
	user, err := ac.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user,
	})
}

// @Summary Logout
// @Description Cerrar sesión (invalidar token del lado del cliente)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	// En una implementación JWT stateless, el logout se maneja del lado del cliente
	// eliminando el token. Aquí solo confirmamos la acción.
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}

// Helper function para obtener user ID desde el contexto
func (ac *AuthController) getUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePublic}
	}
	return userID.(uint), nil
}

// Helper function para obtener user ID desde parámetros
func (ac *AuthController) getUserIDFromParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}