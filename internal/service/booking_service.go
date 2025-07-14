package service

import (
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (s *BookingService) CreateBooking(guestID uuid.UUID, req *models.BookingCreateRequest) (*models.BookingResponse, error) {
	// validate dates
	if req.CheckOut.Before(req.CheckIn) || req.CheckOut.Equal(req.CheckIn) {
		return nil, errors.New("check-out date must be after check-in date")
	}

	if req.CheckIn.Before(time.Now()) {
		return nil, errors.New("check-in date cannot be in the past")
	}

	// get property details
	property, err := s.propertyRepo.GetPropertyByID(req.PropertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		return nil, fmt.Errorf("failed to get property: %w", err)
	}

	// check if property is active
	if property.Status != models.PropertyStatusActive {
		return nil, errors.New("property is not available for booking")
	}

	// check if guest count is within limits
	if req.Guests > property.MaxGuests {
		return nil, fmt.Errorf("number of guests (%d) exceeds property maximum (%d)", req.Guests, property.MaxGuests)
	}

	checkInStr := req.CheckIn.Format("2006-01-02")
	checkOutStr := req.CheckOut.Format("2006-01-02")

	// check for conflicting bookings
	conflictingBookings, err := s.bookingRepo.GetConflictingBookings(req.PropertyID, checkInStr, checkOutStr)
	if err != nil {
		return nil, fmt.Errorf("failed to check for conflicting bookings: %w", err)
	}

	if len(conflictingBookings) > 0 {
		return nil, errors.New("property is not available for the selected dates")
	}

	// check if booking is for at least one night
	nights := int(req.CheckOut.Sub(req.CheckIn).Hours() / 24)
	if nights <= 0 {
		return nil, errors.New("booking must be for at least one night")
	}

	// calculate total price
	totalPrice := float64(nights) * property.PricePerNight

	booking := &models.Booking{
		PropertyID: req.PropertyID,
		GuestID:    guestID,
		CheckIn:    req.CheckIn,
		CheckOut:   req.CheckOut,
		Guests:     req.Guests,
		TotalPrice: totalPrice,
		Currency:   property.Currency,
		Status:     models.BookingStatusPending,
		Notes:      req.Notes,
	}

	err = s.bookingRepo.CreateBooking(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	createdBooking, err := s.bookingRepo.GetBookingByID(booking.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created booking: %w", err)
	}

	return createdBooking.ToResponse(), nil
}

func (s *BookingService) GetBooking(bookingID uuid.UUID, userID uuid.UUID, userRole string) (*models.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if userRole != "admin" && booking.GuestID != userID && booking.Property.HostID != userID {
		return nil, errors.New("unauthorized: you can only view your own bookings")
	}

	return booking.ToResponse(), nil
}

func (s *BookingService) UpdateBooking(bookingID, userID uuid.UUID, userRole string, req *models.BookingUpdateRequest) (*models.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	canUpdate := false
	if userRole == "admin" {
		canUpdate = true
	} else if booking.GuestID == userID {
		canUpdate = true
	} else if booking.Property.HostID == userID {
		canUpdate = true
	}

	if !canUpdate {
		return nil, errors.New("unauthorized: you can only update your own bookings")
	}

	// Update fields based on user role and booking status
	if !req.CheckIn.IsZero() && !req.CheckOut.IsZero() {
		// Only allow date changes if booking is pending and user is the guest
		if booking.Status != models.BookingStatusPending {
			return nil, errors.New("cannot modify dates for confirmed or completed bookings")
		}
		if booking.GuestID != userID && userRole != "admin" {
			return nil, errors.New("only the guest can modify booking dates")
		}

		// Validate new dates
		if req.CheckOut.Before(req.CheckIn) || req.CheckOut.Equal(req.CheckIn) {
			return nil, errors.New("check-out date must be after check-in date")
		}

		if req.CheckIn.Before(time.Now()) {
			return nil, errors.New("check-in date cannot be in the past")
		}

		// Check for conflicts with new dates
		checkInStr := req.CheckIn.Format("2006-01-02")
		checkOutStr := req.CheckOut.Format("2006-01-02")

		conflictingBookings, err := s.bookingRepo.GetConflictingBookings(booking.PropertyID, checkInStr, checkOutStr)
		if err != nil {
			return nil, fmt.Errorf("failed to check for conflicting bookings: %w", err)
		}

		// Filter out the current booking from conflicts
		for _, conflict := range conflictingBookings {
			if conflict.ID != booking.ID {
				return nil, errors.New("property is not available for the selected dates")
			}
		}

		// Update dates and recalculate price
		booking.CheckIn = req.CheckIn
		booking.CheckOut = req.CheckOut
		nights := int(req.CheckOut.Sub(req.CheckIn).Hours() / 24)
		booking.TotalPrice = float64(nights) * booking.Property.PricePerNight
	}

	if req.Guests > 0 {
		// Only allow guest count changes if booking is pending and user is the guest
		if booking.Status != models.BookingStatusPending {
			return nil, errors.New("cannot modify guest count for confirmed or completed bookings")
		}
		if booking.GuestID != userID && userRole != "admin" {
			return nil, errors.New("only the guest can modify guest count")
		}

		// Check guest limits
		if req.Guests > booking.Property.MaxGuests {
			return nil, fmt.Errorf("number of guests (%d) exceeds property maximum (%d)", req.Guests, booking.Property.MaxGuests)
		}
		booking.Guests = req.Guests
	}

	if req.Status != "" {
		// Status changes have specific rules
		switch req.Status {
		case models.BookingStatusConfirmed:
			// Only hosts and admins can confirm bookings
			if booking.Property.HostID != userID && userRole != "admin" {
				return nil, errors.New("only the host can confirm bookings")
			}
			if booking.Status != models.BookingStatusPending {
				return nil, errors.New("only pending bookings can be confirmed")
			}
		case models.BookingStatusCancelled:
			// Guests can cancel pending bookings, hosts and admins can cancel any
			if booking.GuestID == userID && booking.Status == models.BookingStatusPending {
				// Guest cancelling pending booking - allowed
			} else if booking.Property.HostID == userID || userRole == "admin" {
				// Host or admin - allowed
			} else {
				return nil, errors.New("unauthorized to cancel this booking")
			}
		case models.BookingStatusCompleted:
			// Only admins and hosts can mark as completed, and only after check-out date
			if booking.Property.HostID != userID && userRole != "admin" {
				return nil, errors.New("only the host can mark bookings as completed")
			}
			if booking.Status != models.BookingStatusConfirmed {
				return nil, errors.New("only confirmed bookings can be marked as completed")
			}
			if time.Now().Before(booking.CheckOut) {
				return nil, errors.New("booking cannot be completed before check-out date")
			}
		default:
			return nil, errors.New("invalid booking status")
		}
		booking.Status = req.Status
	}

	if req.Notes != "" {
		booking.Notes = req.Notes
	}

	err = s.bookingRepo.UpdateBooking(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	return booking.ToResponse(), nil
}

func (s *BookingService) CancelBooking(bookingID, userID uuid.UUID, userRole string) (*models.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	// Check authorization and booking status
	canCancel := false
	if userRole == "admin" {
		canCancel = true
	} else if booking.GuestID == userID {
		// Guests can cancel their own bookings
		canCancel = true
	} else if booking.Property.HostID == userID {
		// Hosts can cancel bookings for their properties
		canCancel = true
	}

	if !canCancel {
		return nil, errors.New("unauthorized to cancel this booking")
	}

	if booking.Status == models.BookingStatusCancelled {
		return nil, errors.New("booking is already cancelled")
	}

	if booking.Status == models.BookingStatusCompleted {
		return nil, errors.New("cannot cancel completed booking")
	}

	booking.Status = models.BookingStatusCancelled
	err = s.bookingRepo.UpdateBooking(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel booking: %w", err)
	}

	return booking.ToResponse(), nil
}

func (s *BookingService) GetUserBookings(userID uuid.UUID, page, limit int) ([]*models.BookingResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	bookings, err := s.bookingRepo.GetBookingByUserID(userID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}

	responses := make([]*models.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = booking.ToResponse()
	}

	return responses, nil
}

func (s *BookingService) GetPropertyBookings(propertyID, hostID uuid.UUID, page, limit int) ([]*models.BookingResponse, error) {
	// Verify that the user is the host of this property
	property, err := s.propertyRepo.GetPropertyByID(propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		return nil, fmt.Errorf("failed to get property: %w", err)
	}

	if property.HostID != hostID {
		return nil, errors.New("unauthorized: you can only view bookings for your own properties")
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	bookings, err := s.bookingRepo.GetBookingByPropertyID(propertyID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get property bookings: %w", err)
	}

	responses := make([]*models.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = booking.ToResponse()
	}

	return responses, nil
}


