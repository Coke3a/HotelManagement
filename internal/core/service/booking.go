package service

import (
	"context"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type BookingService struct {
	repo port.BookingRepository
}

func NewBookingService(repo port.BookingRepository) *BookingService {
	return &BookingService{
		repo,
	}
}

func (bs *BookingService) CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	if booking.CustomerID == 0 || booking.RatePriceId == 0 || booking.CheckInDate == nil || booking.CheckOutDate == nil || booking.TotalAmount <= 0 {
		return nil, domain.ErrInvalidData
	}

	now := time.Now()
	if booking.BookingDate == nil {
		booking.BookingDate = &now
	}

	// Set initial status to Pending if not provided
	if booking.Status == 0 {
		booking.Status = domain.BookingStatusPending
	}

	createdBooking, err := bs.repo.CreateBooking(ctx, booking)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return createdBooking, nil
}

func (bs *BookingService) GetBooking(ctx context.Context, id uint64) (*domain.Booking, error) {
	booking, err := bs.repo.GetBookingByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return booking, nil
}

func (bs *BookingService) ListBookings(ctx context.Context, skip, limit uint64) ([]domain.Booking, error) {
	bookings, err := bs.repo.ListBookings(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return bookings, nil
}

func (bs *BookingService) ListBookingsWithFilter(ctx context.Context, booking *domain.Booking, skip, limit uint64) ([]domain.Booking, error) {
	// Call the repository's ListBookings method, passing the booking struct, skip, and limit
	bookings, err := bs.repo.ListBookingsWithFilter(ctx, booking, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return bookings, nil
}

func (bs *BookingService) UpdateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	existingBooking, err := bs.repo.GetBookingByID(ctx, booking.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	if booking.CustomerID == existingBooking.CustomerID &&
		booking.RatePriceId == existingBooking.RatePriceId &&
		booking.CheckInDate.Equal(*existingBooking.CheckInDate) &&
		booking.CheckOutDate.Equal(*existingBooking.CheckOutDate) &&
		booking.Status == existingBooking.Status &&
		booking.TotalAmount == existingBooking.TotalAmount {
		return nil, domain.ErrNoUpdatedData
	}

	updatedBooking, err := bs.repo.UpdateBooking(ctx, booking)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return updatedBooking, nil
}

func (bs *BookingService) DeleteBooking(ctx context.Context, id uint64) error {
	_, err := bs.repo.GetBookingByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return bs.repo.DeleteBooking(ctx, id)
}
