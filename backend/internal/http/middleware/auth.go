package middleware

import (
	"net/http"
	"strings"

	"ArtoFino/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	kc *auth.KeycloakClient
}

func NewAuthMiddleware(kc *auth.KeycloakClient) *AuthMiddleware {
	return &AuthMiddleware{kc: kc}
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")

	idToken, err := m.kc.Verifier.Verify(c.Request.Context(), token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	c.Set("claims", claims)
	c.Next()
}
