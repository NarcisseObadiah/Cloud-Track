package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/MicahParks/keyfunc"
)

// Replace with your Zitadel JWKS endpoint
const zitadelJWKSURL = "https://openstack-integration-3vzdfy.us1.zitadel.cloud/oauth/v2/keys"

// Cached JWKS
var jwks *keyfunc.JWKS

func InitJWT() error {
	var err error
	jwks, err = keyfunc.Get(zitadelJWKSURL, keyfunc.Options{})
	return err
}

func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// Extract role (e.g. "admin" or "tenant")
		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing role in token"})
			return
		}

		// Check if user's role matches any of the allowed roles
		if !contains(requiredRoles, role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		// Store role and user ID (optional) in context
		c.Set("role", role)
		c.Set("user_id", claims["sub"])

		c.Next()
	}
}

func contains(list []string, target string) bool {
	for _, val := range list {
		if val == target {
			return true
		}
	}
	return false
}
