package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleGuest UserRole = "guest"
	UserRoleHost  UserRole = "host"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email      string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password   string         `json:"-" gorm:"not null" validate:"required,min=8"`
	FirstName  string         `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	LastName   string         `json:"last_name" gorm:"not null" validate:"required,min=2,max=50"`
	Phone      string         `json:"phone" gorm:"unique"`
	Avatar     string         `json:"avatar"`
	Bio        string         `json:"bio" gorm:"type:text"`
	Role       UserRole       `json:"role" gorm:"type:varchar(20);default:'guest'" validate:"required,oneof=guest host admin"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Properties []Property     `json:"properties,omitempty" gorm:"foreignKey:HostID"`
	Bookings   []Booking      `json:"bookings,omitempty" gorm:"foreignKey:GuestID"`
	Reviews    []Review       `json:"reviews,omitempty" gorm:"foreignKey:ReviewerID"`
}

type UserCreateRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	FirstName string   `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=50"`
	Phone     string   `json:"phone"`
	Role      UserRole `json:"role" validate:"required,oneof=guest host admin"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
