package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"airbnb-clone/internal/service"
	"airbnb-clone/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// creates an authentication middleware
func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// extract token from Bearer header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// validate token
		claims, err := userService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware is similar to AuthMiddleware but doesn't require authentication
func OptionalAuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// extract token from Bearer header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]

		// validate token
		claims, err := userService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// role middleware
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user role"})
			c.Abort()
			return
		}

		// Check if user has required role
		if slices.Contains(roles, roleStr) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// extracts user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	return id, nil
}

// extracts user role from context
func GetUserRole(c *gin.Context) (string, error) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", errors.New("user role not found in context")
	}

	role, ok := userRole.(string)
	if !ok {
		return "", errors.New("invalid user role format")
	}

	return role, nil
}

// extracts user email from context
func GetUserEmail(c *gin.Context) (string, error) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", errors.New("user email not found in context")
	}

	email, ok := userEmail.(string)
	if !ok {
		return "", errors.New("invalid user email format")
	}

	return email, nil
}

// extracts all user claims from context
func GetUserClaims(c *gin.Context) (*utils.Claims, error) {
	userID, err := GetUserID(c)
	if err != nil {
		return nil, err
	}

	userEmail, err := GetUserEmail(c)
	if err != nil {
		return nil, err
	}

	userRole, err := GetUserRole(c)
	if err != nil {
		return nil, err
	}

	return &utils.Claims{
		UserID: userID,
		Email:  userEmail,
		Role:   userRole,
	}, nil
}
