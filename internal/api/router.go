package api

import (
	"airbnb-clone/internal/config"
	"airbnb-clone/internal/middleware"
	"airbnb-clone/internal/service"

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
func NewRouter(services Services, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// global middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

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
		setupAuthRoutes(v1, services.UserService)
		setupPropertyRoutes(v1, services.PropertyService, services.UserService)
		setupBookingRoutes(v1, services.BookingService, services.UserService)
		setupReviewRoutes(v1, services.ReviewService, services.UserService)
	}

	return router
}

// sets up authentication routes
func setupAuthRoutes(rg *gin.RouterGroup, userService *service.UserService) {
	auth := rg.Group("/auth")
	handler := NewUserHandler(userService)

	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.RefreshToken)
}

func setupPropertyRoutes(rg *gin.RouterGroup, propertyService *service.PropertyService, userService *service.UserService) {
	properties := rg.Group("/properties")
	handler := NewPropertyHandler(propertyService)

	// Public routes
	properties.GET("/:id", handler.GetProperty)
	properties.GET("/", handler.ListProperties)

	// Protected routes
	protected := properties.Group("/")
	protected.Use(middleware.AuthMiddleware(userService))
	{
		protected.POST("/", middleware.RequireRole("host", "admin"), handler.CreateProperty)
		protected.PUT("/:id", handler.UpdateProperty)
		protected.DELETE("/:id", handler.DeleteProperty)

	}
}

func setupBookingRoutes(rg *gin.RouterGroup, bookingService *service.BookingService, userService *service.UserService) {

}

func setupReviewRoutes(rg *gin.RouterGroup, reviewService *service.ReviewService, userService *service.UserService) {

}
