package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Coke3a/HotelManagement/internal/adapter/auth/paseto"
	"github.com/Coke3a/HotelManagement/internal/adapter/config"
	"github.com/Coke3a/HotelManagement/internal/adapter/handler/http"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres/repository"
	"github.com/Coke3a/HotelManagement/internal/core/service"
)


func main() {
		// Load environment variables
		config, err := config.New()
		if err != nil {
			slog.Error("Error loading environment variables", "error", err)
			os.Exit(1)
		}
	
		slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)
	
		// Init database
		ctx := context.Background()
		db, err := postgres.Connect(ctx, config.DB)
		if err != nil {
			slog.Error("Error initializing database connection", "error", err)
			os.Exit(1)
		}
		defer db.Close()
	
		slog.Info("Successfully connected to the database", "db", config.DB.Connection)
	
		// Migrate database
		err = db.Migrate()
		if err != nil {
			slog.Error("Error migrating database", "error", err)
			os.Exit(1)
		}
	
		slog.Info("Successfully migrated the database")

		// Init token service
		token, err := paseto.New(config.Token)
		if err != nil {
			slog.Error("Error initializing token service", "error", err)
			os.Exit(1)
		}

		logRepository := repository.NewLogRepository(db)

		userRepository := repository.NewUserRepository(db)
		userService := service.NewUserService(userRepository, logRepository)
		userHandler := http.NewUserHandler(userService)
	

		customerRepository := repository.NewCustomerRepository(db)
		customerService := service.NewCustomerService(customerRepository, logRepository)
		customerHandler := http.NewCustomerHandler(customerService)

		paymentRepository := repository.NewPaymentRepository(db)
		paymentService := service.NewPaymentService(paymentRepository, logRepository)
		paymentHandler := http.NewPaymentHandler(paymentService)

		bookingRepository := repository.NewBookingRepository(db)
		bookingService := service.NewBookingService(bookingRepository, paymentRepository, logRepository)
		bookingHandler := http.NewBookingHandler(bookingService)

		rankRepository := repository.NewRankRepository(db)
		rankService := service.NewRankService(rankRepository, logRepository)
		rankHandler := http.NewRankHandler(rankService)

		ratePriceRepository := repository.NewRatePriceRepository(db)
		ratePriceService := service.NewRatePriceService(ratePriceRepository, logRepository)
		ratePriceHandler := http.NewRatePriceHandler(ratePriceService)

		roomRepository := repository.NewRoomRepository(db)
		roomService := service.NewRoomService(roomRepository, logRepository)
		roomHandler := http.NewRoomHandler(roomService)

		roomTypeRepository := repository.NewRoomTypeRepository(db)
		roomTypeService := service.NewRoomTypeService(roomTypeRepository, logRepository)
		roomTypeHandler := http.NewRoomTypeHandler(roomTypeService)

		customerTypeRepository := repository.NewCustomerTypeRepository(db)
		customerTypeService := service.NewCustomerTypeService(customerTypeRepository, logRepository)
		customerTypeHandler := http.NewCustomerTypeHandler(customerTypeService)

		authService := service.NewAuthService(userRepository, token)
		authHandler := http.NewAuthHandler(authService)
		// Init router
		router, err := http.NewRouter(
			config.HTTP,
			*bookingHandler,
			*customerHandler,
			*paymentHandler,
			*rankHandler,
			*ratePriceHandler,
			*roomHandler,
			*userHandler,
			*authHandler,
			*roomTypeHandler,
			*customerTypeHandler,
			token,
		)
		if err != nil {
			slog.Error("Error initializing router", "error", err)
			os.Exit(1)
		}

		// Start server
		listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
		slog.Info("Starting the HTTP server", "listen_address", listenAddr)
		err = router.Serve(listenAddr)
		if err != nil {
			slog.Error("Error starting the HTTP server", "error", err)
			os.Exit(1)
		}


}