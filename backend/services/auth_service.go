package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"raborimet-crm/backend/config"
	"raborimet-crm/backend/models"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	// Verificar si el email ya existe
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("el email ya está registrado")
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	// Crear el usuario
	user := models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	// Establecer rol por defecto si no se especifica
	if user.Role == "" {
		user.Role = "user"
	}

	// Guardar en la base de datos
	if err := config.DB.Create(&user).Error; err != nil {
		return nil, errors.New("error al crear el usuario")
	}

	return &user, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	// Buscar el usuario por email
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, errors.New("usuario inactivo")
	}

	// Verificar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GenerateToken(user models.User) (string, error) {
	// Crear las claims del token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 días
		"iat":     time.Now().Unix(),
	}

	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token
	tokenString, err := token.SignedString([]byte(s.getJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ChangePassword(userID uint, req models.ChangePasswordRequest) error {
	// Buscar el usuario
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return errors.New("usuario no encontrado")
	}

	// Verificar la contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("contraseña actual incorrecta")
	}

	// Hashear la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error al procesar la nueva contraseña")
	}

	// Actualizar la contraseña
	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		return errors.New("error al actualizar la contraseña")
	}

	return nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	return &user, nil
}

func (s *AuthService) UpdateUserProfile(userID uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Actualizar solo los campos permitidos
	allowedFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if err := config.DB.Model(&user).Updates(filteredUpdates).Error; err != nil {
		return nil, errors.New("error al actualizar el perfil")
	}

	return &user, nil
}

func (s *AuthService) getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "raborimet-crm-secret-key-2024" // Clave por defecto para desarrollo
	}
	return secret
}