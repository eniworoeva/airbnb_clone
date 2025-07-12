package service

import (
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/repository"
	"fmt"

	"github.com/google/uuid"
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

	// Set default currency if not provided
	if property.Currency == "" {
		property.Currency = "USD"
	}

	err := s.propertyRepo.Create(property)
	if err != nil {
		return nil, fmt.Errorf("failed to create property: %w", err)
	}

	// Fetch the created property with host information
	createdProperty, err := s.propertyRepo.GetByID(property.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created property: %w", err)
	}

	return createdProperty.ToResponse(), nil
}