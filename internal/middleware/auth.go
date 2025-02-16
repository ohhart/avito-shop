package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"avito-shop/internal/auth"
	"avito-shop/internal/errors"
	"avito-shop/internal/ports"
)

type AuthMiddleware struct {
	logger *log.Logger
}

func NewAuthMiddleware() ports.AuthMiddleware {
	return &AuthMiddleware{
		logger: log.New(os.Stdout, "AuthMiddleware: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Print("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"errors": errors.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			m.logger.Printf("Invalid token format: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"errors": errors.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]

		claims, err := auth.ParseJWT(tokenString)
		if err != nil {
			m.logger.Printf("JWT parsing error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"errors": errors.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			m.logger.Printf("Invalid claim type for username: %T", claims["username"])
			c.JSON(http.StatusUnauthorized, gin.H{"errors": errors.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
