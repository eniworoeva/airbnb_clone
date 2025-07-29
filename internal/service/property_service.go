package service

import (
	"airbnb-clone/internal/cache"
	"airbnb-clone/internal/logger"
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/repository"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PropertyService struct {
	propertyRepo repository.PropertyRepository
	redisClient  *cache.RedisClient
}

func NewPropertyService(propertyRepo repository.PropertyRepository, redisClient *cache.RedisClient) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
		redisClient:  redisClient,
	}
}

func (s *PropertyService) CreateProperty(hostID uuid.UUID, req *models.PropertyCreateRequest) (*models.PropertyResponse, error) {
	// Create property
	property := &models.Property{
		HostID:        hostID,
		Title:         req.Title,
		Description:   req.Description,
		Type:          req.Type,
		Status:        models.PropertyStatusPending, // Default to pending for approval
		PricePerNight: req.PricePerNight,
		Currency:      req.Currency,
		MaxGuests:     req.MaxGuests,
		Bedrooms:      req.Bedrooms,
		Bathrooms:     req.Bathrooms,
		Address:       req.Address,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Amenities:     req.Amenities,
		Images:        req.Images,
		Rules:         req.Rules,
		CheckInTime:   req.CheckInTime,
		CheckOutTime:  req.CheckOutTime,
	}

	if property.Currency == "" {
		property.Currency = "USD"
	}

	err := s.propertyRepo.Create(property)
	if err != nil {
		logger.Errorf("failed to create property: %v", err)
		return nil, err
	}

	// plans to cache the created property

	createdProperty, err := s.propertyRepo.GetPropertyByID(property.ID)
	if err != nil {
		logger.Errorf("failed to fetch created property: %v", err)
		return nil, err
	}

	return createdProperty.ToResponse(), nil
}

func (s *PropertyService) GetProperty(propertyID uuid.UUID) (*models.PropertyResponse, error) {
	cacheKey := fmt.Sprintf("property:%s", propertyID.String())

	// Try Redis first
	cached, err := s.redisClient.Get(cacheKey)
	if err == nil {
		var property models.Property
		err = json.Unmarshal([]byte(cached), &property)
		if err == nil {
			// Cache hit — return cached property
			return property.ToResponse(), nil
		}
		logger.Warnf("failed to unmarshal cached property: %v", err)
	}

	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		logger.Errorf("failed to get property: %v", err)
		return nil, err
	}

	// Cache the property in Redis
	propertyJSON, err := json.Marshal(property)
	if err != nil {
		logger.Errorf("failed to marshal property for caching: %v", err)
	} else {
		err = s.redisClient.Set(cacheKey, string(propertyJSON), 0) // No expiration
		if err != nil {
			logger.Errorf("failed to cache property: %v", err)
		}
	}

	return property.ToResponse(), nil
}

func (s *PropertyService) UpdateProperty(propertyID, hostID uuid.UUID, req *models.PropertyUpdateRequest) (*models.PropertyResponse, error) {
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		logger.Errorf("failed to get property: %v", err)
		return nil, err
	}

	// check if user is the host of this property
	if property.HostID != hostID {
		return nil, errors.New("unauthorized: you can only update your own properties")
	}

	// Update fields if provided
	if req.Title != "" {
		property.Title = req.Title
	}
	if req.Description != "" {
		property.Description = req.Description
	}
	if req.Type != "" {
		property.Type = req.Type
	}
	if req.Status != "" {
		property.Status = req.Status
	}
	if req.PricePerNight > 0 {
		property.PricePerNight = req.PricePerNight
	}
	if req.Currency != "" {
		property.Currency = req.Currency
	}
	if req.MaxGuests > 0 {
		property.MaxGuests = req.MaxGuests
	}
	if req.Bedrooms >= 0 {
		property.Bedrooms = req.Bedrooms
	}
	if req.Bathrooms > 0 {
		property.Bathrooms = req.Bathrooms
	}
	if req.Address != "" {
		property.Address = req.Address
	}
	if req.City != "" {
		property.City = req.City
	}
	if req.State != "" {
		property.State = req.State
	}
	if req.Country != "" {
		property.Country = req.Country
	}
	if req.ZipCode != "" {
		property.ZipCode = req.ZipCode
	}
	if req.Latitude != 0 {
		property.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		property.Longitude = req.Longitude
	}
	if req.Amenities != nil {
		property.Amenities = req.Amenities
	}
	if req.Images != nil {
		property.Images = req.Images
	}
	if req.Rules != nil {
		property.Rules = req.Rules
	}
	if !req.CheckInTime.IsZero() {
		property.CheckInTime = req.CheckInTime
	}
	if !req.CheckOutTime.IsZero() {
		property.CheckOutTime = req.CheckOutTime
	}

	err = s.propertyRepo.UpdateProperty(property)
	if err != nil {
		logger.Errorf("failed to update property: %v", err)
		return nil, err
	}

	return property.ToResponse(), nil
}

func (s *PropertyService) DeleteProperty(propertyID, hostID uuid.UUID) error {
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("property not found")
		}
		logger.Errorf("failed to get property: %v", err)
		return err
	}

	// Check if user is the host of this property
	if property.HostID != hostID {
		return errors.New("unauthorized: you can only delete your own properties")
	}

	err = s.propertyRepo.DeleteProperty(propertyID)
	if err != nil {
		logger.Errorf("failed to delete property: %v", err)
		return err
	}

	return nil
}

func (s *PropertyService) GetProperties(page, limit int) ([]*models.PropertyResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	properties, err := s.propertyRepo.ListProperties(offset, limit)
	if err != nil {
		logger.Errorf("failed to get properties: %v", err)
		return nil, err
	}

	responses := make([]*models.PropertyResponse, len(properties))
	for i, property := range properties {
		responses[i] = property.ToResponse()
	}

	return responses, nil
}

func (s *PropertyService) SearchProperties(req *models.PropertySearchRequest) (*models.PropertySearchResponse, error) {
	properties, total, err := s.propertyRepo.SearchProperties(req)
	if err != nil {
		logger.Errorf("failed to search properties: %v", err)
		return nil, err
	}

	// Convert to response format
	responses := make([]*models.PropertyResponse, len(properties))
	for i, property := range properties {
		responses[i] = property.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit != 0 {
		totalPages++
	}

	return &models.PropertySearchResponse{
		Properties: responses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *PropertyService) GetPropertiesByHost(hostID uuid.UUID, page, limit int) ([]*models.PropertyResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	properties, err := s.propertyRepo.GetPropertiesByHostID(hostID, offset, limit)
	if err != nil {
		logger.Errorf("failed to get properties by host: %v", err)
		return nil, err
	}

	responses := make([]*models.PropertyResponse, len(properties))
	for i, property := range properties {
		responses[i] = property.ToResponse()
	}

	return responses, nil
}

func (s *PropertyService) CheckAvailability(propertyID uuid.UUID, checkIn, checkOut string) (bool, error) {
	// First check if property exists and is active
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("property not found")
		}
		logger.Errorf("failed to get property: %v", err)
		return false, err
	}

	if property.Status != models.PropertyStatusActive {
		return false, errors.New("property is not available for booking")
	}

	// Check availability
	available, err := s.propertyRepo.CheckAvailability(propertyID, checkIn, checkOut)
	if err != nil {
		logger.Errorf("failed to check availability: %v", err)
		return false, err
	}

	return available, nil
}

func (s *PropertyService) ApproveProperty(propertyID uuid.UUID) (*models.PropertyResponse, error) {
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		logger.Errorf("failed to get property: %v", err)
		return nil, err
	}

	property.Status = models.PropertyStatusActive
	err = s.propertyRepo.UpdateProperty(property)
	if err != nil {
		logger.Errorf("failed to approve property: %v", err)
		return nil, err
	}

	return property.ToResponse(), nil
}
