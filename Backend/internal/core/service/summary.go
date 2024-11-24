package service

import (
	"github.com/gin-gonic/gin"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type SummaryService struct {
	bookingRepo port.BookingRepository
	customerRepo port.CustomerRepository
	roomRepo     port.RoomRepository
	logRepo      port.LogRepository
}

func NewSummaryService(bookingRepo port.BookingRepository, customerRepo port.CustomerRepository, roomRepo port.RoomRepository, logRepo port.LogRepository) *SummaryService {
	return &SummaryService{
		bookingRepo,
		customerRepo,
		roomRepo,
		logRepo,
	}
}

func (ss *SummaryService) GetDashboardSummary(ctx *gin.Context) ([]domain.Booking, []domain.Customer, []domain.Room, error) {

	bookings, _, err := ss.bookingRepo.ListBookings(ctx, 0, 10)
	if err != nil {
		return nil, nil, nil, err
	}

	customers, _, err := ss.customerRepo.ListCustomers(ctx, 0, 10)
	if err != nil {
		return nil, nil, nil, err
	}

	rooms, _, err := ss.roomRepo.ListRooms(ctx, 0, 10)
	if err != nil {
		return nil, nil, nil, err
	}

	return bookings, customers, rooms, nil
}