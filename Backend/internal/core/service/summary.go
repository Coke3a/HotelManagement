package service

import (
	"context"

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

func (ss *SummaryService) GetDashboardSummary(ctx context.Context) ([]domain.Booking, []domain.Customer, []domain.Room, error) {

	bookings, err := ss.bookingRepo.ListBookings(ctx, 0, 10)
	if err != nil {
		return nil, nil, nil, err
	}

	customers, err := ss.customerRepo.ListCustomers(ctx, 0, 10)
	if err != nil {
		return nil, nil, nil, err
	}

	rooms, err := ss.roomRepo.ListRooms(ctx, 0, 10)
	if err != nil {
		return nil, nil, nil, err
	}

	return bookings, customers, rooms, nil
}