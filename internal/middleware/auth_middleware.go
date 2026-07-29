package middleware

import (
	"net/http"
	"strings"

	"visualfinance/internal/pkg/jwt"
	"visualfinance/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware to check for a valid JWT access token
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			c.Abort()
			return
		}

		// Set the UserID in the context so downstream handlers can access it
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}
