package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PropertyType string

const (
	PropertyTypeApartment PropertyType = "apartment"
	PropertyTypeHouse     PropertyType = "house"
	PropertyTypeCondo     PropertyType = "condo"
	PropertyTypeVilla     PropertyType = "villa"
	PropertyTypeCabin     PropertyType = "cabin"
	PropertyTypeStudio    PropertyType = "studio"
)

type PropertyStatus string

const (
	PropertyStatusActive   PropertyStatus = "active"
	PropertyStatusInactive PropertyStatus = "inactive"
	PropertyStatusPending  PropertyStatus = "pending"
)

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

type PropertyUpdateRequest struct {
	Title         string         `json:"title,omitempty" validate:"omitempty,min=10,max=100"`
	Description   string         `json:"description,omitempty" validate:"omitempty,min=50"`
	Type          PropertyType   `json:"type,omitempty" validate:"omitempty,oneof=apartment house condo villa cabin studio"`
	Status        PropertyStatus `json:"status,omitempty" validate:"omitempty,oneof=active inactive pending"`
	PricePerNight float64        `json:"price_per_night,omitempty" validate:"omitempty,min=1"`
	Currency      string         `json:"currency,omitempty"`
	MaxGuests     int            `json:"max_guests,omitempty" validate:"omitempty,min=1,max=20"`
	Bedrooms      int            `json:"bedrooms,omitempty" validate:"omitempty,min=0,max=20"`
	Bathrooms     int            `json:"bathrooms,omitempty" validate:"omitempty,min=1,max=20"`
	Address       string         `json:"address,omitempty"`
	City          string         `json:"city,omitempty"`
	State         string         `json:"state,omitempty"`
	Country       string         `json:"country,omitempty"`
	ZipCode       string         `json:"zip_code,omitempty"`
	Latitude      float64        `json:"latitude,omitempty"`
	Longitude     float64        `json:"longitude,omitempty"`
	Amenities     []string       `json:"amenities,omitempty"`
	Images        []string       `json:"images,omitempty"`
	Rules         []string       `json:"rules,omitempty"`
	CheckInTime   time.Time      `json:"check_in_time"`
	CheckOutTime  time.Time      `json:"check_out_time"`
}

type PropertySearchRequest struct {
	City      string    `json:"city" form:"city"`
	State     string    `json:"state" form:"state"`
	Country   string    `json:"country" form:"country"`
	CheckIn   time.Time `json:"check_in" form:"check_in"`
	CheckOut  time.Time `json:"check_out" form:"check_out"`
	Guests    int       `json:"guests" form:"guests"`
	MinPrice  float64   `json:"min_price" form:"min_price"`
	MaxPrice  float64   `json:"max_price" form:"max_price"`
	Type      string    `json:"type" form:"type"`
	Amenities []string  `json:"amenities" form:"amenities"`
	Page      int       `json:"page" form:"page"`
	Limit     int       `json:"limit" form:"limit"`
}

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

// converts Property to PropertyResponse
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

type PropertySearchResponse struct {
	Properties []*PropertyResponse `json:"properties"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}
