package service

import (
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewService struct {
	reviewRepo  repository.ReviewRepository
	bookingRepo repository.BookingRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, bookingRepo repository.BookingRepository) *ReviewService {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *ReviewService) CreateReview(reviewerID uuid.UUID, req *models.ReviewCreateRequest) (*models.ReviewResponse, error) {
	// Get booking details to validate the review
	booking, err := s.bookingRepo.GetBookingByID(req.BookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if booking.GuestID != reviewerID {
		return nil, errors.New("you can only review bookings you have made")
	}

	if booking.Status != models.BookingStatusCompleted {
		return nil, errors.New("you can only review completed bookings")
	}

	existingReview, err := s.reviewRepo.GetReviewByBookingID(req.BookingID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if existingReview != nil {
		return nil, errors.New("review already exists for this booking")
	}

	review := &models.Review{
		PropertyID: booking.PropertyID,
		BookingID:  req.BookingID,
		ReviewerID: reviewerID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	err = s.reviewRepo.CreateReview(review)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	createdReview, err := s.reviewRepo.GetReviewByID(review.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created review: %w", err)
	}

	return createdReview.ToResponse(), nil
}

func (s *ReviewService) GetReview(reviewID uuid.UUID) (*models.ReviewResponse, error) {
	review, err := s.reviewRepo.GetReviewByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return review.ToResponse(), nil
}

func (s *ReviewService) UpdateReview(reviewID, reviewerID uuid.UUID, req *models.ReviewUpdateRequest) (*models.ReviewResponse, error) {
	review, err := s.reviewRepo.GetReviewByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	if review.ReviewerID != reviewerID {
		return nil, errors.New("you can only update your own reviews")
	}

	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}

	err = s.reviewRepo.UpdateReview(review)
	if err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return review.ToResponse(), nil
}

func (s *ReviewService) DeleteReview(reviewID, reviewerID uuid.UUID, userRole string) error {
	review, err := s.reviewRepo.GetReviewByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return fmt.Errorf("failed to get review: %w", err)
	}

	// Check authorization - reviewer or admin can delete
	if review.ReviewerID != reviewerID && userRole != "admin" {
		return errors.New("you can only delete your own reviews")
	}

	err = s.reviewRepo.DeleteReview(reviewID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

func (s *ReviewService) GetPropertyReviews(propertyID uuid.UUID, page, limit int) ([]*models.ReviewResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	reviews, err := s.reviewRepo.GetReviewsByPropertyID(propertyID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get property reviews: %w", err)
	}

	responses := make([]*models.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = review.ToResponse()
	}

	return responses, nil
}

func (s *ReviewService) GetUserReviews(userID uuid.UUID, page, limit int) ([]*models.ReviewResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	reviews, err := s.reviewRepo.GetReviewsByUserID(userID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}

	responses := make([]*models.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = review.ToResponse()
	}

	return responses, nil
}

// GetAllReviews gets all reviews (admin only)
func (s *ReviewService) GetAllReviews(page, limit int) ([]*models.ReviewResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	reviews, err := s.reviewRepo.ListReviews(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get all reviews: %w", err)
	}

	responses := make([]*models.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = review.ToResponse()
	}

	return responses, nil
}

func (s *ReviewService) GetPropertyAverageRating(propertyID uuid.UUID) (float64, error) {
	avgRating, err := s.reviewRepo.GetAverageRating(propertyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return avgRating, nil
}

