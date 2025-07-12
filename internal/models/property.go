package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	Amenities     pq.StringArray `gorm:"type:text[]" json:"amenities"`
	Images        pq.StringArray `gorm:"type:text[]" json:"images"`
	Rules         pq.StringArray `gorm:"type:text[]" json:"rules"`
	CheckInTime   time.Time      `json:"check_in_time" gorm:"type:time"`
	CheckOutTime  time.Time      `json:"check_out_time" gorm:"type:time"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Host          User           `json:"host,omitempty" gorm:"foreignKey:HostID"`
	Bookings      []Booking      `json:"bookings,omitempty" gorm:"foreignKey:PropertyID"`
	Reviews       []Review       `json:"reviews,omitempty" gorm:"foreignKey:PropertyID"`
}

type PropertyCreateRequest struct {
	Title         string       `json:"title" validate:"required,min=10,max=100"`
	Description   string       `json:"description" validate:"required,min=50"`
	Type          PropertyType `json:"type" validate:"required,oneof=apartment house condo villa cabin studio"`
	PricePerNight float64      `json:"price_per_night" validate:"required,min=1"`
	Currency      string       `json:"currency"`
	MaxGuests     int          `json:"max_guests" validate:"required,min=1,max=20"`
	Bedrooms      int          `json:"bedrooms" validate:"required,min=0,max=20"`
	Bathrooms     int          `json:"bathrooms" validate:"required,min=1,max=20"`
	Address       string       `json:"address" validate:"required"`
	City          string       `json:"city" validate:"required"`
	State         string       `json:"state" validate:"required"`
	Country       string       `json:"country" validate:"required"`
	ZipCode       string       `json:"zip_code" validate:"required"`
	Latitude      float64      `json:"latitude"`
	Longitude     float64      `json:"longitude"`
	Amenities     []string     `json:"amenities"`
	Images        []string     `json:"images"`
	Rules         []string     `json:"rules"`
	CheckInTime   time.Time    `json:"check_in_time"`
	CheckOutTime  time.Time    `json:"check_out_time"`
}

// PropertyResponse represents the response body for property data
type PropertyResponse struct {
	ID            uuid.UUID      `json:"id"`
	HostID        uuid.UUID      `json:"host_id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Type          PropertyType   `json:"type"`
	Status        PropertyStatus `json:"status"`
	PricePerNight float64        `json:"price_per_night"`
	Currency      string         `json:"currency"`
	MaxGuests     int            `json:"max_guests"`
	Bedrooms      int            `json:"bedrooms"`
	Bathrooms     int            `json:"bathrooms"`
	Address       string         `json:"address"`
	City          string         `json:"city"`
	State         string         `json:"state"`
	Country       string         `json:"country"`
	ZipCode       string         `json:"zip_code"`
	Latitude      float64        `json:"latitude"`
	Longitude     float64        `json:"longitude"`
	Amenities     []string       `json:"amenities"`
	Images        []string       `json:"images"`
	Rules         []string       `json:"rules"`
	CheckInTime   time.Time      `json:"check_in_time"`
	CheckOutTime  time.Time      `json:"check_out_time"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Host          *UserResponse  `json:"host,omitempty"`
}

// TableName returns the table name for the Property model
func (Property) TableName() string {
	return "properties"
}

// ToResponse converts Property to PropertyResponse
func (p *Property) ToResponse() *PropertyResponse {
	response := &PropertyResponse{
		ID:            p.ID,
		HostID:        p.HostID,
		Title:         p.Title,
		Description:   p.Description,
		Type:          p.Type,
		Status:        p.Status,
		PricePerNight: p.PricePerNight,
		Currency:      p.Currency,
		MaxGuests:     p.MaxGuests,
		Bedrooms:      p.Bedrooms,
		Bathrooms:     p.Bathrooms,
		Address:       p.Address,
		City:          p.City,
		State:         p.State,
		Country:       p.Country,
		ZipCode:       p.ZipCode,
		Latitude:      p.Latitude,
		Longitude:     p.Longitude,
		Amenities:     p.Amenities,
		Images:        p.Images,
		Rules:         p.Rules,
		CheckInTime:   p.CheckInTime,
		CheckOutTime:  p.CheckOutTime,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if p.Host.ID != uuid.Nil {
		response.Host = p.Host.ToResponse()
	}

	return response
}
