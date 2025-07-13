package service

import (
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PropertyService struct {
	propertyRepo repository.PropertyRepository
}

func NewPropertyService(propertyRepo repository.PropertyRepository) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
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
		return nil, fmt.Errorf("failed to create property: %w", err)
	}

	createdProperty, err := s.propertyRepo.GetPropertyByID(property.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created property: %w", err)
	}

	return createdProperty.ToResponse(), nil
}

func (s *PropertyService) GetProperty(propertyID uuid.UUID) (*models.PropertyResponse, error) {
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		return nil, fmt.Errorf("failed to get property: %w", err)
	}

	return property.ToResponse(), nil
}

func (s *PropertyService) UpdateProperty(propertyID, hostID uuid.UUID, req *models.PropertyUpdateRequest) (*models.PropertyResponse, error) {
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		return nil, fmt.Errorf("failed to get property: %w", err)
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
		return nil, fmt.Errorf("failed to update property: %w", err)
	}

	return property.ToResponse(), nil
}

func (s *PropertyService) DeleteProperty(propertyID, hostID uuid.UUID) error {
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("property not found")
		}
		return fmt.Errorf("failed to get property: %w", err)
	}

	// Check if user is the host of this property
	if property.HostID != hostID {
		return errors.New("unauthorized: you can only delete your own properties")
	}

	err = s.propertyRepo.DeleteProperty(propertyID)
	if err != nil {
		return fmt.Errorf("failed to delete property: %w", err)
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
		return nil, fmt.Errorf("failed to get properties: %w", err)
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
		return nil, fmt.Errorf("failed to search properties: %w", err)
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