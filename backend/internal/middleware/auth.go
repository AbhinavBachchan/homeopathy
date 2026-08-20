package middleware

import (
	"strings"

	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/models"
	"homeopathy-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string      `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// RequireAuth validates the JWT and stashes user_id/role in the gin context.
func RequireAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, 401, "missing or malformed authorization header")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			response.Error(c, 401, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole gate-keeps routes by role, e.g. doctor dashboard, admin panel,
// corporate HR reports.
func RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Error(c, 403, "forbidden")
			c.Abort()
			return
		}
		role := roleVal.(models.Role)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		response.Error(c, 403, "forbidden: insufficient role")
		c.Abort()
	}
}
