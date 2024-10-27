package http

import (
	"github.com/gin-gonic/gin"
	"github.com/Coke3a/HotelManagement/internal/adapter/config"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"net/http"
)

func CORSMiddleware(config *config.HTTP) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.AllowedOrigins)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Access-Control-Expose-Headers")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func ExtractUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenPayload, exists := c.Get("tokenPayload")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		payload, ok := tokenPayload.(*domain.TokenPayload)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token payload"})
			return
		}

		c.Set("user_id", payload.UserID)
		c.Next()
	}
}

func TokenAuthMiddleware(tokenService port.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// Assuming the token is prefixed with "Bearer "
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		payload, err := tokenService.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("tokenPayload", payload)
		c.Next()
	}
}
