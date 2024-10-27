package http

import (
	"github.com/gin-gonic/gin"
	"github.com/Coke3a/HotelManagement/internal/adapter/config"
	// "github.com/Coke3a/HotelManagement/internal/core/domain"
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

// func ExtractUserID() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         tokenPayload, exists := c.Get("tokenPayload")
//         if !exists {
//             c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
//             return
//         }

//         payload, ok := tokenPayload.(*domain.TokenPayload)
//         if !ok {
//             c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token payload"})
//             return
//         }

//         c.Set("userID", payload.UserID)
//         c.Next()
//     }
// }