package repository

import (
	"airbnb-clone/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type propertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) PropertyRepository {
	return &propertyRepository{db: db}
}

func (r *propertyRepository) Create(property *models.Property) error {
	return r.db.Create(property).Error
}

func (r *propertyRepository) GetPropertyByID(id uuid.UUID) (*models.Property, error) {
	var property models.Property
	err := r.db.Preload("Host").Where("id = ?", id).First(&property).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}