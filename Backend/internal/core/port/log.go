package port

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type LogRepository interface {
	CreateLog(ctx context.Context, log *domain.Log) (*domain.Log, error)
}

type LogService interface {
	CreateLog(ctx context.Context, log *domain.Log) (*domain.Log, error)
}
