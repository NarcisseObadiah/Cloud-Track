package auth

import (
	// "fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/MicahParks/keyfunc"
)

const zitadelJWKSURL = "https://openstack-integration-3vzdfy.us1.zitadel.cloud/oauth/v2/keys"

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

		// 🔍 Debugging (optional):
		// fmt.Printf("JWT Claims: %+v\n", claims)

		rolesClaimKey := "urn:zitadel:iam:org:project:roles"

		rawRoles, ok := claims[rolesClaimKey].(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no roles in token"})
			return
		}

		// ✅ Check for required role(s)
		hasRole := false
		for _, required := range requiredRoles {
			if inner, exists := rawRoles[required]; exists {
				if innerMap, ok := inner.(map[string]interface{}); ok && len(innerMap) > 0 {
					hasRole = true
					c.Set("role", required) // Optional
					break
				}
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Set("user_id", claims["sub"])
		c.Next()
	}
}
