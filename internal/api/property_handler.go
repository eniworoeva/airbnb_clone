package api

import (
	"airbnb-clone/internal/middleware"
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	propertyService *service.PropertyService
}

func NewPropertyHandler(propertyService *service.PropertyService) *PropertyHandler {
	return &PropertyHandler{
		propertyService: propertyService,
	}
}

func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req models.PropertyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	property, err := h.propertyService.CreateProperty(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, property)
}