package service

import (
	"fmt"
	"strings"
	"time"
	"log/slog"
	"github.com/gin-gonic/gin"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type DailyBookingSummaryService struct {
	summaryRepo port.DailyBookingSummaryRepository
	bookingRepo port.BookingRepository
	logRepo     port.LogRepository
}

func NewDailyBookingSummaryService(
	summaryRepo port.DailyBookingSummaryRepository,
	bookingRepo port.BookingRepository,
	logRepo port.LogRepository,
) *DailyBookingSummaryService {
	return &DailyBookingSummaryService{
		summaryRepo,
		bookingRepo,
		logRepo,
	}
}

func (dbs *DailyBookingSummaryService) GenerateDailySummary(ctx *gin.Context, date time.Time) (*domain.DailyBookingSummary, error) {
	// Get all bookings for the specified date
	bookings, _, err := dbs.bookingRepo.ListBookingsWithFilter(ctx, &domain.Booking{
		BookingDate: &date,
	}, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("error getting bookings: %w", err)
	}

	// Initialize counters
	summary := &domain.DailyBookingSummary{
		SummaryDate: date,
		Status:      domain.SummaryStatusUnchecked,
	}

	bookingIDs := make([]string, 0, len(bookings))
	
	// Calculate summary statistics
	for _, booking := range bookings {
		summary.TotalBookings++
		summary.TotalAmount += booking.TotalAmount
		bookingIDs = append(bookingIDs, fmt.Sprintf("%d", booking.ID))

		switch booking.Status {
		case domain.BookingStatusPending:
			summary.PendingBookings++
		case domain.BookingStatusConfirmed:
			summary.ConfirmedBookings++
		case domain.BookingStatusCheckedIn:
			summary.CheckedInBookings++
		case domain.BookingStatusCheckedOut:
			summary.CheckedOutBookings++
		case domain.BookingStatusCanceled:
			summary.CanceledBookings++
		case domain.BookingStatusCompleted:
			summary.CompletedBookings++
		}
	}

	summary.BookingIDs = strings.Join(bookingIDs, ",")

	// Create or update the summary
	createdSummary, err := dbs.summaryRepo.CreateDailyBookingSummary(ctx, summary)
	if err != nil {
		return nil, fmt.Errorf("error creating summary: %w", err)
	}

	// Log the action
	userID, exists := ctx.Get("userID")
	if exists {
		log := &domain.Log{
			Action:    "CREATE",
			UserID:    userID.(uint64),
			TableName: "daily_booking_summary",
			RecordID:  uint64(0), // Since summary uses date as primary key
		}
		if _, err := dbs.logRepo.CreateLog(ctx, log); err != nil {
			slog.Error("Error creating log", "error", err)
		}
	}

	return createdSummary, nil
}

func (dbs *DailyBookingSummaryService) UpdateSummaryStatus(ctx *gin.Context, date time.Time, status domain.SummaryStatus) (*domain.DailyBookingSummary, error) {
	summary, err := dbs.summaryRepo.GetDailyBookingSummaryByDate(ctx, date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("error getting summary: %w", err)
	}

	summary.Status = status
	updatedSummary, err := dbs.summaryRepo.UpdateDailyBookingSummary(ctx, summary)
	if err != nil {
		return nil, fmt.Errorf("error updating summary: %w", err)
	}

	// Log the action
	userID, exists := ctx.Get("userID")
	if exists {
		log := &domain.Log{
			Action:    "UPDATE",
			UserID:    userID.(uint64),
			TableName: "daily_booking_summary",
			RecordID:  uint64(0), // Since summary uses date as primary key
		}
		if _, err := dbs.logRepo.CreateLog(ctx, log); err != nil {
			slog.Error("Error creating log", "error", err)
		}
	}

	return updatedSummary, nil
}

func (dbs *DailyBookingSummaryService) GetSummaryByDate(ctx *gin.Context, date time.Time) (*domain.DailyBookingSummary, error) {
	return dbs.summaryRepo.GetDailyBookingSummaryByDate(ctx, date.Format("2006-01-02"))
}

func (dbs *DailyBookingSummaryService) ListSummaries(ctx *gin.Context, skip, limit uint64) ([]domain.DailyBookingSummary, uint64, error) {
	return dbs.summaryRepo.ListDailyBookingSummaries(ctx, skip, limit)
}
