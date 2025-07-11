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
