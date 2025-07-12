package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"airbnb-clone/internal/api"
	"airbnb-clone/internal/config"
	"airbnb-clone/internal/database"
	"airbnb-clone/internal/logger"
	"airbnb-clone/internal/repository"
	"airbnb-clone/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	cfg := config.Load()

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	propertyRepo := repository.NewPropertyRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWT)
	propertyService := service.NewPropertyService(propertyRepo)
	bookingService := service.NewBookingService(bookingRepo, propertyRepo)
	reviewService := service.NewReviewService(reviewRepo, bookingRepo)

	// Initialize router
	router := api.NewRouter(api.Services{
		UserService:     userService,
		PropertyService: propertyService,
		BookingService:  bookingService,
		ReviewService:   reviewService,
	}, cfg)

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Infof("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server gracefully stopped")
	os.Exit(0)
}
