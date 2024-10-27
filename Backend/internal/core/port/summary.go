package port

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type SummaryService interface {
	GetDashboardSummary(ctx context.Context) (*domain.Booking, *domain.Customer, *domain.Room, error)
}
