package service

import (
	"airbnb-clone/internal/repository"
)

type PropertyService struct {
	propertyRepo repository.PropertyRepository
}

func NewPropertyService(propertyRepo repository.PropertyRepository) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
	}
}
