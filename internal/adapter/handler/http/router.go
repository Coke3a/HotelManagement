package http

import (
	"log/slog"
	"strings"

	"github.com/Coke3a/HotelManagement/internal/adapter/config"
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
	bookingHandler BookingHandler,
	customerHandler CustomerHandler,
	paymentHandler PaymentHandler,
	rankHandler RankHandler,
	ratePriceHandler RatePriceHandler,
	roomHandler RoomHandler,
	userHandler UserHandler,
	// authHandler AuthHandler,
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
			user.POST("/", userHandler.CreateUser)
			// user.POST("/login", authHandler.Login)
			user.GET("/", userHandler.ListUsers)
			user.GET("/:id", userHandler.GetUser)
			user.PUT("/:id", userHandler.UpdateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
		}
		conversation := v1.Group("/booking")
		{
			conversation.POST("/", bookingHandler.CreateBooking)
			conversation.GET("/", bookingHandler.ListBookings)
			conversation.GET("/:id", bookingHandler.GetBooking)
			conversation.PUT("/:id", bookingHandler.UpdateBooking)
			conversation.DELETE("/:id", bookingHandler.DeleteBooking)
		}
		customer := v1.Group("/customer")
		{
			customer.POST("/", customerHandler.RegisterCustomer)
			customer.GET("/", customerHandler.ListCustomers)
			customer.GET("/:id", customerHandler.GetCustomer)
			customer.PUT("/:id", customerHandler.UpdateCustomer)
			customer.DELETE("/:id", customerHandler.DeleteCustomer)
		}
		payment := v1.Group("/payment")
		{
			payment.POST("/", paymentHandler.CreatePayment)
			payment.GET("/", paymentHandler.ListPayments)
			payment.GET("/:id", paymentHandler.GetPayment)
			payment.PUT("/:id", paymentHandler.UpdatePayment)
			payment.DELETE("/:id", paymentHandler.DeletePayment)
		}
		rank := v1.Group("/rank")
		{
			rank.POST("/", rankHandler.CreateRank)
			rank.GET("/", rankHandler.ListRanks)
			rank.GET("/:id", rankHandler.GetRank)
			rank.PUT("/:id", rankHandler.UpdateRank)
			rank.DELETE("/:id", rankHandler.DeleteRank)
		}
		ratePrice := v1.Group("/rate_price")
		{
			ratePrice.POST("/", ratePriceHandler.CreateRatePrice)
			ratePrice.GET("/", ratePriceHandler.ListRatePrices)
			ratePrice.GET("/:id", ratePriceHandler.GetRatePrice)
			ratePrice.PUT("/:id", ratePriceHandler.UpdateRatePrice)
			ratePrice.DELETE("/:id", ratePriceHandler.DeleteRatePrice)
		}
		room := v1.Group("/room")
		{
			room.POST("/", roomHandler.CreateRoom)
			room.GET("/", roomHandler.ListRooms)
			room.GET("/:id", roomHandler.GetRoom)
			room.PUT("/:id", roomHandler.UpdateRoom)
			room.DELETE("/:id", roomHandler.DeleteRoom)
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
