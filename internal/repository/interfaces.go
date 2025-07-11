package repository

import (
	"airbnb-clone/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	
}

type PropertyRepository interface {
	
}

type BookingRepository interface {

}

type ReviewRepository interface {
	
}