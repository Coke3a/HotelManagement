package port

import (
	"github.com/gin-gonic/gin"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"time"
)

type DailyBookingSummaryRepository interface {
	CreateDailyBookingSummary(ctx *gin.Context, summary *domain.DailyBookingSummary) (*domain.DailyBookingSummary, error)
	GetDailyBookingSummaryByDate(ctx *gin.Context, date string) (*domain.DailyBookingSummary, error)
	ListDailyBookingSummaries(ctx *gin.Context, skip, limit uint64) ([]domain.DailyBookingSummary, uint64, error)
	UpdateDailyBookingSummary(ctx *gin.Context, summary *domain.DailyBookingSummary) (*domain.DailyBookingSummary, error)
	DeleteDailyBookingSummary(ctx *gin.Context, date string) error
}

type DailyBookingSummaryService interface {
	GenerateDailySummary(ctx *gin.Context, date time.Time) (*domain.DailyBookingSummary, error)
	UpdateSummaryStatus(ctx *gin.Context, date time.Time, status domain.SummaryStatus) (*domain.DailyBookingSummary, error)
	GetSummaryByDate(ctx *gin.Context, date time.Time) (*domain.DailyBookingSummary, error)
	ListSummaries(ctx *gin.Context, skip, limit uint64) ([]domain.DailyBookingSummary, uint64, error)
}
