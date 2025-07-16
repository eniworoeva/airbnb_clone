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
	Property   Property       `json:"property" gorm:"foreignKey:PropertyID"`
	Booking    Booking        `json:"booking" gorm:"foreignKey:BookingID"`
	Reviewer   User           `json:"reviewer" gorm:"foreignKey:ReviewerID"`
}

type ReviewCreateRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Comment   string    `json:"comment" validate:"required,min=10,max=1000"`
}

type ReviewUpdateRequest struct {
	Rating  int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment string `json:"comment,omitempty" validate:"omitempty,min=10,max=1000"`
}

type ReviewResponse struct {
	ID         uuid.UUID         `json:"id"`
	PropertyID uuid.UUID         `json:"property_id"`
	BookingID  uuid.UUID         `json:"booking_id"`
	ReviewerID uuid.UUID         `json:"reviewer_id"`
	Rating     int               `json:"rating"`
	Comment    string            `json:"comment"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Property   *PropertyResponse `json:"property,omitempty"`
	Reviewer   *UserResponse     `json:"reviewer,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}

func (r *Review) ToResponse() *ReviewResponse {
	response := &ReviewResponse{
		ID:         r.ID,
		PropertyID: r.PropertyID,
		BookingID:  r.BookingID,
		ReviewerID: r.ReviewerID,
		Rating:     r.Rating,
		Comment:    r.Comment,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}

	if r.Property.ID != uuid.Nil {
		response.Property = r.Property.ToResponse()
	}

	if r.Reviewer.ID != uuid.Nil {
		response.Reviewer = r.Reviewer.ToResponse()
	}

	return response
}

type PropertyRatingResponse struct {
	PropertyID    uuid.UUID `json:"property_id"`
	AverageRating float64   `json:"average_rating"`
	ReviewCount   int       `json:"review_count"`
}
