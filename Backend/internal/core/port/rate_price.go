package port

import (
	"github.com/gin-gonic/gin"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type RatePriceRepository interface {
	CreateRatePrice(ctx *gin.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	GetRatePriceByID(ctx *gin.Context, id uint64) (*domain.RatePrice, error)
	ListRatePrices(ctx *gin.Context, skip, limit uint64) ([]domain.RatePrice, error)
	UpdateRatePrice(ctx *gin.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	DeleteRatePrice(ctx *gin.Context, id uint64) error
	GetRatePricesByRoomTypeId(ctx *gin.Context, roomTypeID uint64) ([]domain.RatePrice, error)
	GetRatePricesByRoomId(ctx *gin.Context, roomID uint64) ([]domain.RatePrice, error)
}

type RatePriceService interface {
	CreateRatePrice(ctx *gin.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	GetRatePrice(ctx *gin.Context, id uint64) (*domain.RatePrice, error)
	ListRatePrices(ctx *gin.Context, skip, limit uint64) ([]domain.RatePrice, error)
	UpdateRatePrice(ctx *gin.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error)
	DeleteRatePrice(ctx *gin.Context, id uint64) error
	GetRatePricesByRoomTypeId(ctx *gin.Context, roomTypeID uint64) ([]domain.RatePrice, error)
	GetRatePricesByRoomId(ctx *gin.Context, roomID uint64) ([]domain.RatePrice, error)
}