package middleware

import (
	"net/http"

	"ArtoFino/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func RequireRole(role string, clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(map[string]interface{})

		if !auth.HasRole(claims, role, clientID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden: missing role " + role,
			})
			return
		}

		c.Next()
	}
}
