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

		token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// Zitadel roles claim key
		rolesClaimKey := "urn:zitadel:iam:org:project:roles"

		rolesMap, ok := claims[rolesClaimKey].(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no roles in token"})
			return
		}

		// Check if user has one of the required roles
		hasRole := false
		for _, role := range requiredRoles {
			if _, ok := rolesMap[role]; ok {
				hasRole = true
				c.Set("role", role) // optional: store role in context
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Set("user_id", claims["sub"]) // optional
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
