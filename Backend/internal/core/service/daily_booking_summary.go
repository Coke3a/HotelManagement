package service

import (
	"fmt"
	"time"
	"log/slog"
	"strconv"
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
	// Get all created bookings for the specified date
	createdBookings, createdCount, err := dbs.bookingRepo.ListBookingsWithFilter(ctx, &domain.Booking{
		CreatedAt: &date,
	}, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("error getting created bookings: %w", err)
	}

	// Get all updated bookings for the specified date
	updatedBookings, updatedCount, err := dbs.bookingRepo.ListBookingsWithFilter(ctx, &domain.Booking{
		UpdatedAt: &date,
	}, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("error getting updated bookings: %w", err)
	}

	slog.Info("Retrieved bookings",
		"date", date.Format("2006-01-02"),
		"created_count", createdCount,
		"updated_count", updatedCount)

	// Initialize arrays for different booking types
	var (
		createdIDs    []uint64
		completedIDs  []uint64
		canceledIDs   []uint64
		totalAmount   float64
	)

	// Process created bookings
	for _, booking := range createdBookings {
		createdIDs = append(createdIDs, booking.ID)
	}

	// Process updated bookings
	for _, booking := range updatedBookings {
		// Check if the booking was already counted in created bookings
		if !contains(createdIDs, booking.ID) {
			if booking.Status == domain.BookingStatusCompleted {
				totalAmount += booking.TotalAmount
				completedIDs = append(completedIDs, booking.ID)
			} else if booking.Status == domain.BookingStatusCanceled {
				canceledIDs = append(canceledIDs, booking.ID)
			}
		}
	}

	// Create summary object
	summary := &domain.DailyBookingSummary{
		SummaryDate:       date,
		CreatedBookings:   domain.FormatBookingIDs(createdIDs),
		CompletedBookings: domain.FormatBookingIDs(completedIDs),
		CanceledBookings:  domain.FormatBookingIDs(canceledIDs),
		TotalAmount:       totalAmount,
		Status:            domain.SummaryStatusUnchecked,
	}

	// Create or update the summary
	createdSummary, err := dbs.summaryRepo.CreateDailyBookingSummary(ctx, summary)
	if err != nil {
		return nil, fmt.Errorf("error creating summary: %w", err)
	}

	// Log the action
	userID, exists := ctx.Get("userID")
	if exists {
		dateStr := date.Format("20060102") // YYYYMMDD format
		recordID, err := strconv.ParseUint(dateStr, 10, 64)
		if err != nil {
			slog.Error("Error parsing date for log record ID", "error", err)
			recordID = 0 // Use 0 as fallback
		}
		
		log := &domain.Log{
			Action:    "CREATE",
			UserID:    userID.(uint64),
			TableName: "daily_booking_summary",
			RecordID:  recordID,
		}
		if _, err := dbs.logRepo.CreateLog(ctx, log); err != nil {
			slog.Error("Error creating log", "error", err)
		}
	}

	return createdSummary, nil
}

// Helper function to check if a slice contains a value
func contains(slice []uint64, value uint64) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
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
