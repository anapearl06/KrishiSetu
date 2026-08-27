package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/raaj2493/KrishiSetu/internal/auth"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			response.Error(
				c,
				401,
				"UNAUTHORIZED",
				"missing authorization token",
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(
				c,
				401,
				"UNAUTHORIZED",
				"invalid authorization header",
			)
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(
			parts[1],
			secret,
		)
		if err != nil {
			response.Error(
				c,
				401,
				"UNAUTHORIZED",
				"invalid or expired token",
			)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)

		c.Next()
	}
}