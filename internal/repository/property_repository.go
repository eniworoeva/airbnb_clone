package repository

import (
	"airbnb-clone/internal/models"
	"fmt"
	"strings"

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

func (r *propertyRepository) UpdateProperty(property *models.Property) error {
	return r.db.Save(property).Error
}

func (r *propertyRepository) DeleteProperty(id uuid.UUID) error {
	return r.db.Delete(&models.Property{}, id).Error
}

func (r *propertyRepository) ListProperties(offset, limit int) ([]*models.Property, error) {
	var properties []*models.Property
	err := r.db.Preload("Host").Where("status = ?", "active").Offset(offset).Limit(limit).Find(&properties).Error
	return properties, err
}

func (r *propertyRepository) SearchProperties(req *models.PropertySearchRequest) ([]*models.Property, int64, error) {
	var properties []*models.Property
	var total int64

	// Build the base query
	query := r.db.Model(&models.Property{}).Preload("Host")
	countQuery := r.db.Model(&models.Property{})

	// Apply filters
	whereClause := "status = 'active'"
	args := []interface{}{}

	if req.City != "" {
		whereClause += " AND LOWER(city) LIKE LOWER(?)"
		args = append(args, "%"+req.City+"%")
	}

	if req.State != "" {
		whereClause += " AND LOWER(state) LIKE LOWER(?)"
		args = append(args, "%"+req.State+"%")
	}

	if req.Country != "" {
		whereClause += " AND LOWER(country) LIKE LOWER(?)"
		args = append(args, "%"+req.Country+"%")
	}

	if req.Guests > 0 {
		whereClause += " AND max_guests >= ?"
		args = append(args, req.Guests)
	}

	if req.MinPrice > 0 {
		whereClause += " AND price_per_night >= ?"
		args = append(args, req.MinPrice)
	}

	if req.MaxPrice > 0 {
		whereClause += " AND price_per_night <= ?"
		args = append(args, req.MaxPrice)
	}

	if req.Type != "" {
		whereClause += " AND type = ?"
		args = append(args, req.Type)
	}

	// Handle amenities search using raw SQL for array operations
	if len(req.Amenities) > 0 {
		amenitiesPlaceholders := make([]string, len(req.Amenities))
		for i, amenity := range req.Amenities {
			amenitiesPlaceholders[i] = "?"
			args = append(args, amenity)
		}
		whereClause += fmt.Sprintf(" AND amenities ?& ARRAY[%s]", strings.Join(amenitiesPlaceholders, ","))
	}

	// Apply availability filter if check-in and check-out dates are provided
	if !req.CheckIn.IsZero() && !req.CheckOut.IsZero() {
		subQuery := `
			NOT EXISTS (
				SELECT 1 FROM bookings 
				WHERE property_id = properties.id 
				AND status IN ('confirmed', 'pending')
				AND NOT (check_out <= ? OR check_in >= ?)
			)
		`
		whereClause += " AND " + subQuery
		args = append(args, req.CheckIn, req.CheckOut)
	}

	// Apply where clause to both queries
	query = query.Where(whereClause, args...)
	countQuery = countQuery.Where(whereClause, args...)

	// Get total count
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	// Execute the main query
	err = query.Offset(offset).Limit(req.Limit).Find(&properties).Error
	if err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *propertyRepository) GetPropertiesByHostID(hostID uuid.UUID, offset, limit int) ([]*models.Property, error) {
	var properties []*models.Property
	err := r.db.Where("host_id = ?", hostID).Offset(offset).Limit(limit).Find(&properties).Error
	return properties, err
}
