package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PropertyID uuid.UUID      `json:"property_id" gorm:"type:uuid;not null"`
	BookingID  uuid.UUID      `json:"booking_id" gorm:"type:uuid;not null"`
	ReviewerID uuid.UUID      `json:"reviewer_id" gorm:"type:uuid;not null"`
	Rating     int            `json:"rating" gorm:"not null" validate:"required,min=1,max=5"`
	Comment    string         `json:"comment" gorm:"type:text" validate:"required,min=10,max=1000"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Property   Property       `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
	Booking    Booking        `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	Reviewer   User           `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
}
