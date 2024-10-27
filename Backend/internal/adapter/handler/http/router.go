package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/config"
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
	authHandler AuthHandler,
	roomTypeHandler RoomTypeHandler,
	customerTypeHandler CustomerTypeHandler,
) (*Router, error) {
	// Disable debug mode in production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(CORSMiddleware(config))

	//
	// START LOG FOR DEBUGGING
	//
	router.Use(func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse the body as JSON
		var bodyJSON interface{}
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &bodyJSON); err != nil {
				bodyJSON = string(bodyBytes) // If not JSON, use the raw string
			}
		}

		// Log the request
		slog.Info("Incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_addr", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"body", bodyJSON)

		c.Next()
	})
	//
	// END LOG FOR DEBUGGING
	//

	router.Use(sloggin.New(slog.Default()), gin.Recovery())

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	// v1.Use(ExtractUserID()) // Add this line
	{

		user := v1.Group("/users")
		{
			user.POST("/", userHandler.CreateUser)
			user.POST("/login", authHandler.Login)
			user.GET("/:id", userHandler.GetUser)
			user.GET("/", userHandler.ListUsers)
			user.PUT("/", userHandler.UpdateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
		}
		booking := v1.Group("/booking")
		{
			booking.POST("/", bookingHandler.CreateBookingAndPayment)
			booking.GET("/", bookingHandler.ListBookingCustomerPayments)
			booking.GET("/:id", bookingHandler.GetBooking)
			booking.PUT("/", bookingHandler.UpdateBooking)
			booking.DELETE("/:id", bookingHandler.DeleteBooking)
			booking.GET("/:id/details", bookingHandler.GetBookingCustomerPayment)
		}
		customer := v1.Group("/customers")
		{
			customer.POST("/", customerHandler.CreateCustomer)
			customer.GET("/", customerHandler.ListCustomers)
			customer.GET("/:id", customerHandler.GetCustomer)
			customer.PUT("/", customerHandler.UpdateCustomer)
			customer.DELETE("/:id", customerHandler.DeleteCustomer)
		}
		payment := v1.Group("/payments")
		{
			payment.POST("/", paymentHandler.CreatePayment)
			payment.GET("/", paymentHandler.ListPayments)
			payment.GET("/:id", paymentHandler.GetPayment)
			payment.PUT("/", paymentHandler.UpdatePayment)
			payment.DELETE("/:id", paymentHandler.DeletePayment)
		}
		rank := v1.Group("/ranks")
		{
			rank.POST("/", rankHandler.CreateRank)
			rank.GET("/", rankHandler.ListRanks)
			rank.GET("/:id", rankHandler.GetRank)
			rank.PUT("/", rankHandler.UpdateRank)
			rank.DELETE("/:id", rankHandler.DeleteRank)
		}
		ratePrice := v1.Group("/rate_prices")
		{
			ratePrice.POST("/", ratePriceHandler.CreateRatePrice)
			ratePrice.GET("/", ratePriceHandler.ListRatePrices)
			ratePrice.GET("/:id", ratePriceHandler.GetRatePrice)
			ratePrice.PUT("/", ratePriceHandler.UpdateRatePrice)
			ratePrice.DELETE("/:id", ratePriceHandler.DeleteRatePrice)
			ratePrice.GET("/by-room-type/:room_type_id", ratePriceHandler.GetRatePricesByRoomTypeId)
		}
		room := v1.Group("/rooms")
		{
			room.POST("/", roomHandler.CreateRoom)
			room.GET("/", roomHandler.ListRooms)
			room.GET("/:id", roomHandler.GetRoom)
			room.PUT("/", roomHandler.UpdateRoom)
			room.DELETE("/:id", roomHandler.DeleteRoom)
			room.GET("/available", roomHandler.GetAvailableRooms)
			room.GET("/with-room-type", roomHandler.ListRoomsWithRoomType)
		}
		roomTypeRoutes := v1.Group("/room-types")
		{
			roomTypeRoutes.POST("/", roomTypeHandler.CreateRoomType)
			roomTypeRoutes.GET("/:id", roomTypeHandler.GetRoomType)
			roomTypeRoutes.GET("/", roomTypeHandler.ListRoomTypes)
			roomTypeRoutes.PUT("/", roomTypeHandler.UpdateRoomType)
			roomTypeRoutes.DELETE("/:id", roomTypeHandler.DeleteRoomType)
		}
		customerTypeRoutes := v1.Group("/customer-types")
		{
			customerTypeRoutes.POST("/", customerTypeHandler.CreateCustomerType)
			customerTypeRoutes.GET("/:id", customerTypeHandler.GetCustomerType)
			customerTypeRoutes.GET("/", customerTypeHandler.ListCustomerTypes)
			customerTypeRoutes.PUT("/", customerTypeHandler.UpdateCustomerType)
			customerTypeRoutes.DELETE("/:id", customerTypeHandler.DeleteCustomerType)
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
