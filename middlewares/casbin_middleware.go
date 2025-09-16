package middlewares

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el rol del contexto (previamente guardado por AuthMiddleware)
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Role not found in context",
			})
			return
		}

		obj := c.FullPath()
		act := c.Request.Method

		// Debug info (opcional, remover en producción)
		fmt.Printf("Casbin Check - Role: %s, Object: %s, Action: %s\n", role, obj, act)

		// Verificar permisos con Casbin
		ok, err := e.Enforce(role, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Permission check failed",
			})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"details": fmt.Sprintf("Role '%s' cannot '%s' on '%s'", role, act, obj),
			})
			return
		}

		c.Next()
	}
}
