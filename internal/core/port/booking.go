package port

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetBookingByID(ctx context.Context, id uint64) (*domain.Booking, error)
	ListBookings(ctx context.Context, skip, limit uint64) ([]domain.Booking, error)
	UpdateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	DeleteBooking(ctx context.Context, id uint64) error
}

type BookingService interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetBooking(ctx context.Context, id uint64) (*domain.Booking, error)
	ListBookings(ctx context.Context, skip, limit uint64) ([]domain.Booking, error)
	UpdateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	DeleteBooking(ctx context.Context, id uint64) error
}
