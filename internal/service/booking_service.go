package service

import (
	"airbnb-clone/internal/repository"
)

type BookingService struct {
	bookingRepo  repository.BookingRepository
	propertyRepo repository.PropertyRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, propertyRepo repository.PropertyRepository) *BookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		propertyRepo: propertyRepo,
	}
}
