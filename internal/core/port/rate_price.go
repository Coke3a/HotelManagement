package port

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type RatePriceRepository interface {
	CreateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	GetRatePriceByID(ctx context.Context, id uint64) (*domain.RatePrice, error)
	ListRatePrices(ctx context.Context, skip, limit uint64) ([]domain.RatePrice, error)
	UpdateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	DeleteRatePrice(ctx context.Context, id uint64) error
}

type RatePriceService interface {
	RegisterRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	GetRatePrice(ctx context.Context, id uint64) (*domain.RatePrice, error)
	ListRatePrices(ctx context.Context, skip, limit uint64) ([]domain.RatePrice, error)
	UpdateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	DeleteRatePrice(ctx context.Context, id uint64) error
}
