package api

import (
	"airbnb-clone/internal/config"
	"airbnb-clone/internal/middleware"
	"airbnb-clone/internal/cache"
	"airbnb-clone/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

// holds all service dependencies
type Services struct {
	UserService     *service.UserService
	PropertyService *service.PropertyService
	BookingService  *service.BookingService
	ReviewService   *service.ReviewService
}

// creates and configures the main router
func NewRouter(services Services, cfg *config.Config, redisClient *cache.RedisClient) *gin.Engine {
	router := gin.New()

	// global middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Global rate limiting middleware
	rateLimitConfig := middleware.RateLimiterConfig{
		RequestsPerMinute: cfg.RateLimit.DefaultRequestsPerMinute,
		WindowSize:        time.Minute,
		KeyPrefix:         "rate_limit",
	}
	router.Use(middleware.RateLimiterMiddleware(redisClient, rateLimitConfig))

	// health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "airbnb-clone",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		setupAuthRoutes(v1, services.UserService, redisClient, cfg)
		setupPropertyRoutes(v1, services.PropertyService, services.UserService, redisClient, cfg)
		setupBookingRoutes(v1, services.BookingService, services.UserService, redisClient, cfg)
		setupReviewRoutes(v1, services.ReviewService, services.UserService, redisClient, cfg)
	}

	return router
}

// sets up authentication routes
func setupAuthRoutes(rg *gin.RouterGroup, userService *service.UserService, redisClient *cache.RedisClient, cfg *config.Config) {
	auth := rg.Group("/auth")
	handler := NewUserHandler(userService)

	// Apply stricter rate limiting for auth endpoints
	authRateLimit := middleware.CreateRateLimiterForEndpoint(redisClient, cfg.RateLimit.AuthRequestsPerMinute, "auth")
	auth.Use(authRateLimit)

	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.RefreshToken)
}

func setupPropertyRoutes(rg *gin.RouterGroup, propertyService *service.PropertyService, userService *service.UserService, redisClient *cache.RedisClient, cfg *config.Config) {
	properties := rg.Group("/properties")
	handler := NewPropertyHandler(propertyService)

	// Public routes with moderate rate limiting
	properties.GET("/", handler.ListProperties)
	properties.GET("/search", middleware.CreateRateLimiterForEndpoint(redisClient, cfg.RateLimit.SearchRequestsPerMinute, "search"), handler.SearchProperties)
	properties.GET("/:id", handler.GetProperty)
	properties.GET("/:id/availability", handler.CheckAvailability)

	// Protected routes
	protected := properties.Group("/")
	protected.Use(middleware.AuthMiddleware(userService))
	{
		protected.POST("/", middleware.RequireRole("host", "admin"), handler.CreateProperty)
		protected.PUT("/:id", handler.UpdateProperty)
		protected.DELETE("/:id", handler.DeleteProperty)
		protected.GET("/my", handler.GetMyProperties)

		// Admin only routes
		admin := protected.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("/:id/approve", handler.ApproveProperty)
		}

	}
}

func setupBookingRoutes(rg *gin.RouterGroup, bookingService *service.BookingService, userService *service.UserService, redisClient *cache.RedisClient, cfg *config.Config) {
	bookings := rg.Group("/bookings")
	bookings.Use(middleware.AuthMiddleware(userService))
	handler := NewBookingHandler(bookingService)

	// Apply stricter rate limiting for booking creation
	bookings.POST("/", middleware.CreateRateLimiterForEndpoint(redisClient, cfg.RateLimit.BookingRequestsPerMinute, "create_booking"), handler.CreateBooking)
	bookings.GET("/:id", handler.GetBooking)
	bookings.PUT("/:id", handler.UpdateBooking)
	bookings.POST("/:id/cancel", handler.CancelBooking)
	bookings.GET("/my", handler.GetMyBookings)
	bookings.GET("/property/:property_id", handler.GetPropertyBookings)

	// Admin only routes
	admin := bookings.Group("/")
	admin.Use(middleware.RequireRole("admin"))

}

func setupReviewRoutes(rg *gin.RouterGroup, reviewService *service.ReviewService, userService *service.UserService, redisClient *cache.RedisClient, cfg *config.Config) {
	reviews := rg.Group("/reviews")
	handler := NewReviewHandler(reviewService)

	protected := reviews.Group("/")
	protected.Use(middleware.AuthMiddleware(userService))
	{
		// Apply moderate rate limiting for review creation
		protected.POST("/", middleware.CreateRateLimiterForEndpoint(redisClient, cfg.RateLimit.ReviewRequestsPerMinute, "create_review"), handler.CreateReview)
		protected.GET("/:id", handler.GetReview)
		protected.PUT("/:id", handler.UpdateReview)
		protected.DELETE("/:id", handler.DeleteReview)
	}
}
