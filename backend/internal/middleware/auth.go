package middleware

import (
	"strings"

	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/models"
	"homeopathy-platform/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string      `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// RequireAuth validates the JWT and stashes user_id/role in the gin context.
func RequireAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return response.Error(c, 401, "missing or malformed authorization header")
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return response.Error(c, 401, "invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

// RequireRole gate-keeps routes by role, e.g. doctor dashboard, admin panel,
// corporate HR reports.
func RequireRole(roles ...models.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(models.Role) // or string, whatever type claims.Role actually is
		if !ok {
			return response.Error(c, 401, "role missing from context")
		}

		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return response.Error(c, 403, "forbidden: insufficient role")

	}
}
