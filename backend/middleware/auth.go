package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log de debugging
		fmt.Printf("[AUTH DEBUG] Petición: %s %s\n", c.Request.Method, c.Request.URL.Path)
		
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("[AUTH DEBUG] Authorization header: '%s'\n", authHeader)
		
		if authHeader == "" {
			fmt.Printf("[AUTH DEBUG] No authorization header found, returning 401\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización requerido"})
			c.Abort()
			return
		}

		// Verificar formato Bearer
		tokenParts := strings.Split(authHeader, " ")
		fmt.Printf("[AUTH DEBUG] Token parts: %v\n", tokenParts)
		
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			fmt.Printf("[AUTH DEBUG] Invalid token format, returning 401\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		fmt.Printf("[AUTH DEBUG] Token string: %s...\n", tokenString[:min(len(tokenString), 20)])

		// Validar el token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(getJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			fmt.Printf("[AUTH DEBUG] Token validation failed: %v, returning 401\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Verificar que el usuario existe y está activo
		var user models.User
		fmt.Printf("[AUTH DEBUG] Looking for user ID: %d\n", claims.UserID)
		
		if err := config.DB.First(&user, claims.UserID).Error; err != nil {
			fmt.Printf("[AUTH DEBUG] User not found: %v, returning 401\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
			c.Abort()
			return
		}

		if !user.IsActive {
			fmt.Printf("[AUTH DEBUG] User inactive, returning 401\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario inactivo"})
			c.Abort()
			return
		}

		fmt.Printf("[AUTH DEBUG] Authentication successful for user: %s\n", user.Email)
		
		// Guardar información del usuario en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user", user)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
			c.Abort()
			return
		}

		// Verificar si el rol del usuario está en la lista de roles permitidos
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Permisos insuficientes"})
		c.Abort()
	}
}

func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := tokenParts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(getJWTSecret()), nil
		})

		if err == nil && token.Valid {
			var user models.User
			if err := config.DB.First(&user, claims.UserID).Error; err == nil && user.IsActive {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
				c.Set("user", user)
			}
		}

		c.Next()
	}
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "raborimet-crm-secret-key-2024" // Clave por defecto para desarrollo
	}
	return secret
}