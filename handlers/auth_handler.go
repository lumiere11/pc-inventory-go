package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lumiere11/pc-inventory-go/models"
	"github.com/lumiere11/pc-inventory-go/requests"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret_key")

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}
func (h *AuthHandler) verifyPasswords(password, password_confirmation string) bool {
	return password == password_confirmation
}

// Registro de usuario
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := context.Background()
	var req requests.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if !h.verifyPasswords(req.Password, req.PasswordConfirmation) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	var user models.User
	user.Email = req.Email
	user.Password = string(hash)
	user.Role = "normal_user"
	result := h.DB.WithContext(ctx).Create(&user)
	if result.Error != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{}, "message": "User created"})
}

// Login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.User
	ctx := context.Background()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	var user models.User

	result := h.DB.WithContext(ctx).Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": result.Error.Error(),
		})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	exp := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
