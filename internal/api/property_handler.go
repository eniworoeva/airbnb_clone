package api

import (
	"airbnb-clone/internal/middleware"
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *PropertyHandler) GetProperty(c *gin.Context) {
	propertyIDStr := c.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	property, err := h.propertyService.GetProperty(propertyID)
	if err != nil {
		if err.Error() == "property not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	propertyIDStr := c.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	var req models.PropertyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	property, err := h.propertyService.UpdateProperty(propertyID, userID, &req)
	if err != nil {
		if err.Error() == "property not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized: you can only update your own properties" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	propertyIDStr := c.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	err = h.propertyService.DeleteProperty(propertyID, userID)
	if err != nil {
		if err.Error() == "property not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized: you can only delete your own properties" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Property deleted successfully"})
}

func (h *PropertyHandler) ListProperties(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	properties, err := h.propertyService.GetProperties(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"properties": properties,
		"page":       page,
		"limit":      limit,
	})
}

func (h *PropertyHandler) SearchProperties(c *gin.Context) {
	var req models.PropertySearchRequest

	req.City = c.Query("city")
	req.State = c.Query("state")
	req.Country = c.Query("country")
	req.Type = c.Query("type")
	req.Amenities = c.QueryArray("amenities")

	if checkIn := c.Query("check_in"); checkIn != "" {
		parsed, err := time.Parse("2006-01-02", checkIn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_in format. Use YYYY-MM-DD"})
			return
		}
		req.CheckIn = parsed
	}

	if checkOut := c.Query("check_out"); checkOut != "" {
		parsed, err := time.Parse("2006-01-02", checkOut)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_out format. Use YYYY-MM-DD"})
			return
		}
		req.CheckOut = parsed
	}

	if guests := c.Query("guests"); guests != "" {
		num, err := strconv.Atoi(guests)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guests value"})
			return
		}
		req.Guests = num
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		price, err := strconv.ParseFloat(minPrice, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_price"})
			return
		}
		req.MinPrice = price
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		price, err := strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_price"})
			return
		}
		req.MaxPrice = price
	}

	req.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	req.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.propertyService.SearchProperties(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PropertyHandler) GetMyProperties(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	properties, err := h.propertyService.GetPropertiesByHost(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"properties": properties,
		"page":       page,
		"limit":      limit,
	})
}