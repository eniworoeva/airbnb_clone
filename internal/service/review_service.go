package service

import (
	"airbnb-clone/internal/repository"
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
