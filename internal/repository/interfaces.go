package repository

import (
	"airbnb-clone/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	
}

type PropertyRepository interface {
	Create(property *models.Property) error
	GetPropertyByID(id uuid.UUID) (*models.Property, error) 
	UpdateProperty(property *models.Property) error
	DeleteProperty(id uuid.UUID) error 
	ListProperties(offset, limit int) ([]*models.Property, error)
	SearchProperties(req *models.PropertySearchRequest) ([]*models.Property, int64, error) 
}

type BookingRepository interface {

}

type ReviewRepository interface {
	
}