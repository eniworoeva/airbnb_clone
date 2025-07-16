package repository

import (
	"airbnb-clone/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	
}

type PropertyRepository interface {
	Create(property *models.Property) error
	GetPropertyByID(id uuid.UUID) (*models.Property, error) 
	UpdateProperty(property *models.Property) error
	DeleteProperty(id uuid.UUID) error 
	ListProperties(offset, limit int) ([]*models.Property, error)
	SearchProperties(req *models.PropertySearchRequest) ([]*models.Property, int64, error) 
	GetPropertiesByHostID(hostID uuid.UUID, offset, limit int) ([]*models.Property, error)
	CheckAvailability(propertyID uuid.UUID, checkIn, checkOut string) (bool, error)
}

type BookingRepository interface {
	GetConflictingBookings(propertyID uuid.UUID, checkIn, checkOut string) ([]*models.Booking, error)
	CreateBooking(booking *models.Booking) error 
	GetBookingByID(id uuid.UUID) (*models.Booking, error)
	GetBookingByUserID(userID uuid.UUID, offset, limit int) ([]*models.Booking, error)
	GetBookingByPropertyID(propertyID uuid.UUID, offset, limit int) ([]*models.Booking, error)
	UpdateBooking(booking *models.Booking) error
	DeleteBooking(id uuid.UUID) error
}

type ReviewRepository interface {
	CreateReview(review *models.Review) error
	GetReviewByID(id uuid.UUID) (*models.Review, error)
	GetReviewsByPropertyID(propertyID uuid.UUID, offset, limit int) ([]*models.Review, error)
	GetReviewsByUserID(userID uuid.UUID, offset, limit int) ([]*models.Review, error)
	GetReviewByBookingID(bookingID uuid.UUID) (*models.Review, error)
	UpdateReview(review *models.Review) error
	DeleteReview(id uuid.UUID) error
	ListReviews(offset, limit int) ([]*models.Review, error)
	GetAverageRating(propertyID uuid.UUID) (float64, error)
}