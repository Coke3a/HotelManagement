package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type SummaryService interface {
	GetDashboardSummary(ctx *gin.Context) (*domain.Booking, *domain.Customer, *domain.Room, error)
}
