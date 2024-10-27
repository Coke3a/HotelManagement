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
	GetRatePricesByRoomTypeId(ctx context.Context, roomTypeID uint64) ([]domain.RatePrice, error)
	GetRatePricesByRoomId(ctx context.Context, roomID uint64) ([]domain.RatePrice, error)
}

type RatePriceService interface {
	CreateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	GetRatePrice(ctx context.Context, id uint64) (*domain.RatePrice, error)
	ListRatePrices(ctx context.Context, skip, limit uint64) ([]domain.RatePrice, error)
	UpdateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	DeleteRatePrice(ctx context.Context, id uint64) error
	GetRatePricesByRoomTypeId(ctx context.Context, roomTypeID uint64) ([]domain.RatePrice, error)
	GetRatePricesByRoomId(ctx context.Context, roomID uint64) ([]domain.RatePrice, error)
}