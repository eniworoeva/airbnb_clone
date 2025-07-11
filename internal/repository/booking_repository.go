package repository

import (

	"gorm.io/gorm"
)

// implements BookingRepository interface
type bookingRepository struct {
	db *gorm.DB
}

// creates a new booking repository
func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}
