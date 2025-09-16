package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener el header "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Falta el header de autorización"})
			return
		}

		// 2. Validar que el formato sea "Bearer <token>" y extraer el token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header de autorización mal formado"})
			return
		}
		tokenString := parts[1]

		// 3. Parsear y validar el token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Asegurarse de que el algoritmo de firma es el esperado
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo de firma inesperado: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		// 4. Manejar errores de validación
		if err != nil || !token.Valid {
			// Se pueden revisar errores específicos como token expirado
			if err == jwt.ErrTokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "El token ha expirado"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// 5. Si el token es válido, guardar información en el contexto
		// para que los handlers posteriores puedan usarla.
		c.Set("email", claims.Email)

		// 6. Continuar con el siguiente handler en la cadena.
		c.Next()
	}
}
