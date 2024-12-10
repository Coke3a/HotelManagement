package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type BookingRepository interface {
	CreateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error)
	GetBookingByID(ctx *gin.Context, id uint64) (*domain.Booking, error)
	ListBookings(ctx *gin.Context, skip, limit uint64) ([]domain.Booking, uint64, error)
	ListBookingsWithFilter(ctx *gin.Context, booking *domain.Booking, skip, limit uint64) ([]domain.Booking, uint64, error)
	UpdateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error)
	DeleteBooking(ctx *gin.Context, id uint64) error
	GetBookingCustomerPayment(ctx *gin.Context, id uint64) (*domain.BookingCustomerPayment, error)
	ListBookingCustomerPayments(ctx *gin.Context, skip, limit uint64) ([]domain.BookingCustomerPayment, uint64, error)
	ListBookingCustomerPaymentsWithFilter(ctx *gin.Context, bookingCustomerPayment *domain.BookingCustomerPayment, skip, limit uint64) ([]domain.BookingCustomerPayment, uint64, error)
}

type BookingService interface {
	CreateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error)
	GetBooking(ctx *gin.Context, id uint64) (*domain.Booking, error)
	ListBookings(ctx *gin.Context, skip, limit uint64) ([]domain.Booking, uint64, error)
	ListBookingsWithFilter(ctx *gin.Context, booking *domain.Booking, skip, limit uint64) ([]domain.Booking, uint64,error)
	CreateBookingAndPayment(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error)
	UpdateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error)
	DeleteBooking(ctx *gin.Context, id uint64) error
	GetBookingCustomerPayment(ctx *gin.Context, id uint64) (*domain.BookingCustomerPayment, error)
	ListBookingCustomerPayments(ctx *gin.Context, skip, limit uint64) ([]domain.BookingCustomerPayment, uint64, error)
	ListBookingCustomerPaymentsWithFilter(ctx *gin.Context, bookingCustomerPayment *domain.BookingCustomerPayment, skip, limit uint64) ([]domain.BookingCustomerPayment, uint64, error)
}
