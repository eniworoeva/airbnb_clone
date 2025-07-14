package repository

import (
	"airbnb-clone/internal/models"

	"github.com/google/uuid"
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

func (r *bookingRepository) GetConflictingBookings(propertyID uuid.UUID, checkIn, checkOut string) ([]*models.Booking, error) {
	var bookings []*models.Booking
	
	query := `
		SELECT * FROM bookings 
		WHERE property_id = ? 
		AND status IN ('confirmed', 'pending')
		AND NOT (check_out <= ? OR check_in >= ?)
	`
	
	err := r.db.Raw(query, propertyID, checkIn, checkOut).Scan(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) CreateBooking(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) GetBookingByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Property").Preload("Guest").Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetBookingByUserID(userID uuid.UUID, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.Preload("Property").Preload("Guest").
		Where("guest_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetBookingByPropertyID(propertyID uuid.UUID, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.Preload("Property").Preload("Guest").
		Where("property_id = ?", propertyID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) UpdateBooking(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) DeleteBooking(id uuid.UUID) error {
	return r.db.Delete(&models.Booking{}, id).Error
}
