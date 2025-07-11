package api

import (
	"airbnb-clone/internal/service"
)

type PropertyHandler struct {
	propertyService *service.PropertyService
}

func NewPropertyHandler(propertyService *service.PropertyService) *PropertyHandler {
	return &PropertyHandler{
		propertyService: propertyService,
	}
}
