package api

import (
	"airbnb-clone/internal/service"

)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}
