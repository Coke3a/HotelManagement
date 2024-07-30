package http

import (
	"log/slog"
	"strings"

	"github.com/Coke3a/TalkPenguin/internal/adapter/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new HTTP router
func NewRouter(
	config *config.HTTP,
	// token port.TokenService,
	// userHandler UserHandler,
	// authHandler AuthHandler,
	messageHandler MessageHandler,
	conversationHandler ConversationHandler,
	userHandler UserHandler,
) (*Router, error) {
	// Disable debug mode in production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	// // Custom validators
	// v, ok := binding.Validator.Engine().(*validator.Validate)
	// if ok {
	// 	if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
	// 		return nil, err
	// 	}

	// 	if err := v.RegisterValidation("payment_type", paymentTypeValidator); err != nil {
	// 		return nil, err
	// 	}

	// }

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{

		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			// user.POST("/login", authHandler.Login)
			user.GET("/", userHandler.ListUsers)
			user.GET("/:id", userHandler.GetUser)
			user.PUT("/:id", userHandler.UpdateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
		}
		conversation := v1.Group("/conversations")
		{
			conversation.POST("/create", conversationHandler.CreateConversation)
			conversation.POST("/create_with_message", conversationHandler.CreateConversationWithMessage)
		}
		message := v1.Group("/messages")
		{
			message.POST("/exchange", messageHandler.ExchangingMessage)
			message.GET("/all/:user_id/:conversation_id", messageHandler.GetAllMessages)
		}
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
