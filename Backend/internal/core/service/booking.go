package service

import (
	"log/slog"
	"time"
	"fmt"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
)

type BookingService struct {
	repo       port.BookingRepository
	paymentRepo port.PaymentRepository
	logRepo     port.LogRepository
}

func NewBookingService(repo port.BookingRepository, paymentRepo port.PaymentRepository, logRepo port.LogRepository) *BookingService {
	return &BookingService{
		repo,
		paymentRepo,
		logRepo,
	}
}

func (bs *BookingService) CreateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error) {
	if booking.CustomerID == 0 || booking.RatePriceId == 0 || booking.CheckInDate == nil || booking.CheckOutDate == nil || booking.TotalAmount <= 0 {
		return nil, domain.ErrInvalidData
	}

	// Set initial status to Pending if not provided
	if booking.Status == 0 {
		booking.Status = domain.BookingStatusUncheckIn
	}

	createdBooking, err := bs.repo.CreateBooking(ctx, booking)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  createdBooking.ID,
		Action:    "CREATE",
		UserID:    userID.(uint64),
		TableName: "bookings",
	}
	_, err = bs.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return createdBooking, nil
}

func (bs *BookingService) CreateBookingAndPayment(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error) {
	if booking.CustomerID == 0 || booking.RatePriceId == 0 || booking.CheckInDate == nil || booking.CheckOutDate == nil || booking.TotalAmount <= 0 {
		return nil, domain.ErrInvalidData
	}

	now := time.Now()
	// Set initial status to Pending if not provided
	if booking.Status == 0 {
		booking.Status = domain.BookingStatusUncheckIn
	}

	// Create the booking
	createdBooking, err := bs.repo.CreateBooking(ctx, booking)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Create the payment record
	payment := &domain.Payment{
		BookingID:     createdBooking.ID,
		Amount:        booking.TotalAmount,
		PaymentMethod: domain.PaymentMethodNotSpecified,
		PaymentDate:   &now,
		Status:        domain.PaymentStatusUnpaid,
	}

	// Create the payment
	_, err = bs.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		// Rollback the booking creation if payment creation fails
		_ = bs.repo.DeleteBooking(ctx, createdBooking.ID)
		return nil, domain.ErrInternal
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  createdBooking.ID,
		Action:    "CREATE",
		UserID:    userID.(uint64),
		TableName: "bookings",
	}
	_, err = bs.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return createdBooking, nil
}

func (bs *BookingService) GetBooking(ctx *gin.Context, id uint64) (*domain.Booking, error) {
	booking, err := bs.repo.GetBookingByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return booking, nil
}

func (bs *BookingService) ListBookings(ctx *gin.Context, skip, limit uint64) ([]domain.Booking, uint64, error) {
	bookings, totalCount, err := bs.repo.ListBookings(ctx, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}

	return bookings, totalCount, nil
}

func (bs *BookingService) ListBookingsWithFilter(ctx *gin.Context, booking *domain.Booking, skip, limit uint64) ([]domain.Booking, uint64, error) {
	bookings, totalCount, err := bs.repo.ListBookingsWithFilter(ctx, booking, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}

	return bookings, totalCount, nil
}

func (bs *BookingService) UpdateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error) {
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
		booking.RoomID == existingBooking.RoomID &&
		booking.RoomTypeID == existingBooking.RoomTypeID &&
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

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  booking.ID,
		Action:    "UPDATE",
		UserID:    userID.(uint64),
		TableName: "bookings",
	}
	_, err = bs.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return updatedBooking, nil
}

func (bs *BookingService) DeleteBooking(ctx *gin.Context, id uint64) error {
	_, err := bs.repo.GetBookingByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		return domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  id,
		Action:    "DELETE",
		UserID:    userID.(uint64),
		TableName: "bookings",
	}
	_, err = bs.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return bs.repo.DeleteBooking(ctx, id)
}

func (bs *BookingService) GetBookingCustomerPayment(ctx *gin.Context, id uint64) (*domain.BookingCustomerPayment, error) {
	bookingCustomerPayment, err := bs.repo.GetBookingCustomerPayment(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return bookingCustomerPayment, nil
}

func (bs *BookingService) ListBookingCustomerPayments(ctx *gin.Context, skip, limit uint64) ([]domain.BookingCustomerPayment, uint64, error) {
	bookingCustomerPayments, totalCount, err := bs.repo.ListBookingCustomerPayments(ctx, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}

	return bookingCustomerPayments, totalCount, nil
}

func (bs *BookingService) ListBookingCustomerPaymentsWithFilter(ctx *gin.Context, bookingCustomerPayment *domain.BookingCustomerPayment, skip, limit uint64) ([]domain.BookingCustomerPayment, uint64, error) {
	bookingCustomerPayments, totalCount, err := bs.repo.ListBookingCustomerPaymentsWithFilter(ctx, bookingCustomerPayment, skip, limit)
	if err != nil {
		fmt.Println("Error in ListBookingCustomerPaymentsWithFilter", err)
		return nil, 0, domain.ErrInternal
	}

	return bookingCustomerPayments, totalCount, nil
}