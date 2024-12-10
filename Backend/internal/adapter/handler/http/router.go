package http

import (

	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/config"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/Coke3a/HotelManagement/internal/core/port"
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
	logHandler LogHandler,
	dailyBookingSummaryHandler DailyBookingSummaryHandler,
	tokenService port.TokenService,
) (*Router, error) {
	router := SetupRouter(config, tokenService)

	v1 := router.Group("/v1")
	{
		// Public routes (no authentication required)
		user := v1.Group("/users")
		{
			user.POST("/login", authHandler.Login)    // Login
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(TokenAuthMiddleware(tokenService))
		protected.Use(ExtractUserID())
		{
			// User routes
			userProtected := protected.Group("/users")
			{
				userProtected.POST("/", userHandler.CreateUser)
				userProtected.GET("/:id", userHandler.GetUser)
				userProtected.GET("/", userHandler.ListUsers)
				userProtected.PUT("/", userHandler.UpdateUser)
				userProtected.DELETE("/:id", userHandler.DeleteUser)
			}

			// Other protected routes
			booking := protected.Group("/booking")
			{
				booking.POST("/", bookingHandler.CreateBookingAndPayment)
				booking.GET("/", bookingHandler.ListBookingCustomerPaymentsWithFilter)
				// booking.GET("/", bookingHandler.ListBookingsWithFilter)
				booking.GET("/:id", bookingHandler.GetBooking)
				booking.PUT("/", bookingHandler.UpdateBooking)
				booking.DELETE("/:id", bookingHandler.DeleteBooking)
				booking.GET("/:id/details", bookingHandler.GetBookingCustomerPayment)
			}
			customer := protected.Group("/customers")
			{
				customer.POST("/", customerHandler.CreateCustomer)
				customer.GET("/", customerHandler.ListCustomers)
				customer.GET("/:id", customerHandler.GetCustomer)
				customer.PUT("/", customerHandler.UpdateCustomer)
				customer.DELETE("/:id", customerHandler.DeleteCustomer)
			}
			payment := protected.Group("/payments")
			{
				payment.POST("/", paymentHandler.CreatePayment)
				payment.GET("/", paymentHandler.ListPayments)
				payment.GET("/:id", paymentHandler.GetPayment)
				payment.PUT("/", paymentHandler.UpdatePayment)
				payment.DELETE("/:id", paymentHandler.DeletePayment)
			}
			rank := protected.Group("/ranks")
			{
				rank.POST("/", rankHandler.CreateRank)
				rank.GET("/", rankHandler.ListRanks)
				rank.GET("/:id", rankHandler.GetRank)
				rank.PUT("/", rankHandler.UpdateRank)
				rank.DELETE("/:id", rankHandler.DeleteRank)
			}
			ratePrice := protected.Group("/rate_prices")
			{
				ratePrice.POST("/", ratePriceHandler.CreateRatePrice)
				ratePrice.GET("/", ratePriceHandler.ListRatePrices)
				ratePrice.GET("/:id", ratePriceHandler.GetRatePrice)
				ratePrice.PUT("/", ratePriceHandler.UpdateRatePrice)
				ratePrice.DELETE("/:id", ratePriceHandler.DeleteRatePrice)
				ratePrice.GET("/by-room-type/:room_type_id", ratePriceHandler.GetRatePricesByRoomTypeId)
			}
			room := protected.Group("/rooms")
			{
				room.POST("/", roomHandler.CreateRoom)
				room.GET("/", roomHandler.ListRooms)
				room.GET("/:id", roomHandler.GetRoom)
				room.PUT("/", roomHandler.UpdateRoom)
				room.DELETE("/:id", roomHandler.DeleteRoom)
				room.GET("/available", roomHandler.GetAvailableRooms)
				room.GET("/with-room-type", roomHandler.ListRoomsWithRoomType)
			}
			roomTypeRoutes := protected.Group("/room-types")
			{
				roomTypeRoutes.POST("/", roomTypeHandler.CreateRoomType)
				roomTypeRoutes.GET("/:id", roomTypeHandler.GetRoomType)
				roomTypeRoutes.GET("/", roomTypeHandler.ListRoomTypes)
				roomTypeRoutes.PUT("/", roomTypeHandler.UpdateRoomType)
				roomTypeRoutes.DELETE("/:id", roomTypeHandler.DeleteRoomType)
			}
			customerTypeRoutes := protected.Group("/customer-types")
			{
				customerTypeRoutes.POST("/", customerTypeHandler.CreateCustomerType)
				customerTypeRoutes.GET("/:id", customerTypeHandler.GetCustomerType)
				customerTypeRoutes.GET("/", customerTypeHandler.ListCustomerTypes)
				customerTypeRoutes.PUT("/", customerTypeHandler.UpdateCustomerType)
				customerTypeRoutes.DELETE("/:id", customerTypeHandler.DeleteCustomerType)
			}
			dailySummary := protected.Group("/daily-summary")
			{
				dailySummary.POST("/generate", dailyBookingSummaryHandler.GenerateDailySummary)
				dailySummary.PUT("/status", dailyBookingSummaryHandler.UpdateSummaryStatus)
				dailySummary.GET("/", dailyBookingSummaryHandler.GetSummaryByDate)
				dailySummary.GET("/list", dailyBookingSummaryHandler.ListSummaries)
			}
			log := protected.Group("/logs")
			{
				log.GET("/", logHandler.GetLogs)
			}
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

func SetupRouter(config *config.HTTP, tokenService port.TokenService) *gin.Engine {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(CORSMiddleware(config))

	router.Use(sloggin.New(slog.Default()), gin.Recovery())

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
