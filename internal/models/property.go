package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PropertyType represents the type of property
type PropertyType string

const (
	PropertyTypeApartment PropertyType = "apartment"
	PropertyTypeHouse     PropertyType = "house"
	PropertyTypeCondo     PropertyType = "condo"
	PropertyTypeVilla     PropertyType = "villa"
	PropertyTypeCabin     PropertyType = "cabin"
	PropertyTypeStudio    PropertyType = "studio"
)

// PropertyStatus represents the status of a property
type PropertyStatus string

const (
	PropertyStatusActive   PropertyStatus = "active"
	PropertyStatusInactive PropertyStatus = "inactive"
	PropertyStatusPending  PropertyStatus = "pending"
)

// Property represents the property entity
type Property struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	HostID        uuid.UUID      `json:"host_id" gorm:"type:uuid;not null"`
	Title         string         `json:"title" gorm:"not null" validate:"required,min=10,max=100"`
	Description   string         `json:"description" gorm:"type:text" validate:"required,min=50"`
	Type          PropertyType   `json:"type" gorm:"type:varchar(20);not null" validate:"required,oneof=apartment house condo villa cabin studio"`
	Status        PropertyStatus `json:"status" gorm:"type:varchar(20);default:'pending'" validate:"required,oneof=active inactive pending"`
	PricePerNight float64        `json:"price_per_night" gorm:"not null" validate:"required,min=1"`
	Currency      string         `json:"currency" gorm:"default:'USD'"`
	MaxGuests     int            `json:"max_guests" gorm:"not null" validate:"required,min=1,max=20"`
	Bedrooms      int            `json:"bedrooms" gorm:"not null" validate:"required,min=0,max=20"`
	Bathrooms     int            `json:"bathrooms" gorm:"not null" validate:"required,min=1,max=20"`
	Address       string         `json:"address" gorm:"not null" validate:"required"`
	City          string         `json:"city" gorm:"not null" validate:"required"`
	State         string         `json:"state" gorm:"not null" validate:"required"`
	Country       string         `json:"country" gorm:"not null" validate:"required"`
	ZipCode       string         `json:"zip_code" gorm:"not null" validate:"required"`
	Latitude      float64        `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude     float64        `json:"longitude" gorm:"type:decimal(11,8)"`
	Amenities     []string       `json:"amenities" gorm:"type:text[];serializer:json"`
	Images        []string       `json:"images" gorm:"type:text[];serializer:json"`
	Rules         []string       `json:"rules" gorm:"type:text[];serializer:json"`
	CheckInTime   time.Time      `json:"check_in_time" gorm:"type:time"`
	CheckOutTime  time.Time      `json:"check_out_time" gorm:"type:time"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Host          User           `json:"host,omitempty" gorm:"foreignKey:HostID"`
	Bookings      []Booking      `json:"bookings,omitempty" gorm:"foreignKey:PropertyID"`
	Reviews       []Review       `json:"reviews,omitempty" gorm:"foreignKey:PropertyID"`
}
