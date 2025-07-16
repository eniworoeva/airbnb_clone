package repository

import (
	"airbnb-clone/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) CreateReview(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetReviewByID(id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("Property").Preload("Reviewer").Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetReviewsByPropertyID(propertyID uuid.UUID, offset, limit int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.Preload("Reviewer").
		Where("property_id = ?", propertyID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetReviewsByUserID(userID uuid.UUID, offset, limit int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.Preload("Property").
		Where("reviewer_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetReviewByBookingID(bookingID uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("Property").Preload("Reviewer").
		Where("booking_id = ?", bookingID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) UpdateReview(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) DeleteReview(id uuid.UUID) error {
	return r.db.Delete(&models.Review{}, id).Error
}

func (r *reviewRepository) ListReviews(offset, limit int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.Preload("Property").Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetAverageRating(propertyID uuid.UUID) (float64, error) {
	var avgRating float64

	query := `
		SELECT COALESCE(AVG(rating), 0) as avg_rating 
		FROM reviews 
		WHERE property_id = ? AND deleted_at IS NULL
	`

	err := r.db.Raw(query, propertyID).Scan(&avgRating).Error
	return avgRating, err
}
