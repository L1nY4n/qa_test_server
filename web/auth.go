package web

import (
	"net/http"
	"strings"

	"qa_test_server/manager"

	"github.com/gin-gonic/gin"
)

const authClaimsKey = "authClaims"

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			fail(c, http.StatusUnauthorized, "missing authorization token")
			c.Abort()
			return
		}

		claims, err := manager.UserManagerGlobal.ParseToken(token)
		if err != nil {
			fail(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(authClaimsKey, claims)
		c.Next()
	}
}

func requireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[strings.ToLower(strings.TrimSpace(role))] = struct{}{}
	}

	return func(c *gin.Context) {
		claims, ok := currentClaims(c)
		if !ok {
			fail(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		if _, exists := allowed[strings.ToLower(claims.Role)]; !exists {
			fail(c, http.StatusForbidden, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

func currentClaims(c *gin.Context) (manager.TokenClaims, bool) {
	v, ok := c.Get(authClaimsKey)
	if !ok {
		return manager.TokenClaims{}, false
	}
	claims, ok := v.(manager.TokenClaims)
	return claims, ok
}

func bearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	if len(header) < 8 {
		return ""
	}
	if strings.ToLower(header[:7]) != "bearer " {
		return ""
	}
	return strings.TrimSpace(header[7:])
}
