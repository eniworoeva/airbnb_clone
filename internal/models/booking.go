package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PropertyID uuid.UUID      `json:"property_id" gorm:"type:uuid;not null"`
	GuestID    uuid.UUID      `json:"guest_id" gorm:"type:uuid;not null"`
	CheckIn    time.Time      `json:"check_in" gorm:"not null" validate:"required"`
	CheckOut   time.Time      `json:"check_out" gorm:"not null" validate:"required"`
	Guests     int            `json:"guests" gorm:"not null" validate:"required,min=1"`
	TotalPrice float64        `json:"total_price" gorm:"not null" validate:"required,min=0"`
	Currency   string         `json:"currency" gorm:"default:'USD'"`
	Status     BookingStatus  `json:"status" gorm:"type:varchar(20);default:'pending'" validate:"required,oneof=pending confirmed cancelled completed"`
	Notes      string         `json:"notes" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Property   Property       `json:"property,omitzero" gorm:"foreignKey:PropertyID"`
	Guest      User           `json:"guest,omitzero" gorm:"foreignKey:GuestID"`
}

type BookingCreateRequest struct {
	PropertyID uuid.UUID `json:"property_id" validate:"required"`
	CheckIn    time.Time `json:"check_in" validate:"required"`
	CheckOut   time.Time `json:"check_out" validate:"required"`
	Guests     int       `json:"guests" validate:"required,min=1"`
	Notes      string    `json:"notes"`
}

type BookingUpdateRequest struct {
	CheckIn  time.Time     `json:"check_in,omitempty"`
	CheckOut time.Time     `json:"check_out,omitempty"`
	Guests   int           `json:"guests,omitempty" validate:"omitempty,min=1"`
	Status   BookingStatus `json:"status,omitempty" validate:"omitempty,oneof=pending confirmed cancelled completed"`
	Notes    string        `json:"notes,omitempty"`
}

type BookingResponse struct {
	ID         uuid.UUID         `json:"id"`
	PropertyID uuid.UUID         `json:"property_id"`
	GuestID    uuid.UUID         `json:"guest_id"`
	CheckIn    time.Time         `json:"check_in"`
	CheckOut   time.Time         `json:"check_out"`
	Guests     int               `json:"guests"`
	TotalPrice float64           `json:"total_price"`
	Currency   string            `json:"currency"`
	Status     BookingStatus     `json:"status"`
	Notes      string            `json:"notes"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Property   *PropertyResponse `json:"property,omitempty"`
	Guest      *UserResponse     `json:"guest,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}

// ToResponse converts Booking to BookingResponse
func (b *Booking) ToResponse() *BookingResponse {
	response := &BookingResponse{
		ID:         b.ID,
		PropertyID: b.PropertyID,
		GuestID:    b.GuestID,
		CheckIn:    b.CheckIn,
		CheckOut:   b.CheckOut,
		Guests:     b.Guests,
		TotalPrice: b.TotalPrice,
		Currency:   b.Currency,
		Status:     b.Status,
		Notes:      b.Notes,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}

	if b.Property.ID != uuid.Nil {
		response.Property = b.Property.ToResponse()
	}

	if b.Guest.ID != uuid.Nil {
		response.Guest = b.Guest.ToResponse()
	}

	return response
}