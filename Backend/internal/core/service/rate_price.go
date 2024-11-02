package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type RatePriceService struct {
	repo    port.RatePriceRepository
	logRepo port.LogRepository
}

func NewRatePriceService(repo port.RatePriceRepository, logRepo port.LogRepository) *RatePriceService {
	return &RatePriceService{
		repo,
		logRepo,
	}
}

func (rps *RatePriceService) CreateRatePrice(ctx *gin.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error) {
	if ratePrice.Name == "" || ratePrice.PricePerNight < 0 || ratePrice.RoomTypeID == 0 {
		return nil, domain.ErrInvalidData
	}

	createdRatePrice, err := rps.repo.CreateRatePrice(ctx, ratePrice)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  createdRatePrice.ID,
		Action:    "CREATE",
		UserID:    userID.(uint64),
		TableName: "rate_prices",
	}
	_, err = rps.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return createdRatePrice, nil
}

func (rps *RatePriceService) GetRatePrice(ctx *gin.Context, id uint64) (*domain.RatePrice, error) {
	ratePrice, err := rps.repo.GetRatePriceByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return ratePrice, nil
}

func (rps *RatePriceService) ListRatePrices(ctx *gin.Context, skip, limit uint64) ([]domain.RatePrice, error) {
	ratePrices, err := rps.repo.ListRatePrices(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return ratePrices, nil
}

func (rps *RatePriceService) UpdateRatePrice(ctx *gin.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error) {
	existingRatePrice, err := rps.repo.GetRatePriceByID(ctx, ratePrice.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	emptyData := ratePrice.Name == "" &&
		ratePrice.Description == "" &&
		ratePrice.PricePerNight == 0 &&
		ratePrice.RoomTypeID == 0

	sameData := existingRatePrice.Name == ratePrice.Name &&
		existingRatePrice.Description == ratePrice.Description &&
		existingRatePrice.PricePerNight == ratePrice.PricePerNight &&
		existingRatePrice.RoomTypeID == ratePrice.RoomTypeID

	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	// Update timestamp
	now := time.Now()
	ratePrice.UpdatedAt = &now

	updatedRatePrice, err := rps.repo.UpdateRatePrice(ctx, ratePrice)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  ratePrice.ID,
		Action:    "UPDATE",
		UserID:    userID.(uint64),
		TableName: "rate_prices",
	}
	_, err = rps.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return updatedRatePrice, nil
}

func (rps *RatePriceService) DeleteRatePrice(ctx *gin.Context, id uint64) error {
	_, err := rps.repo.GetRatePriceByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		return domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  id,
		Action:    "DELETE",
		UserID:    userID.(uint64),
		TableName: "rate_prices",
	}
	_, err = rps.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return rps.repo.DeleteRatePrice(ctx, id)
}

func (rps *RatePriceService) GetRatePricesByRoomTypeId(ctx *gin.Context, roomTypeID uint64) ([]domain.RatePrice, error) {
	ratePrices, err := rps.repo.GetRatePricesByRoomTypeId(ctx, roomTypeID)
	if err != nil {
		return nil, domain.ErrInternal
	}

	if len(ratePrices) == 0 {
		return nil, domain.ErrDataNotFound
	}

	return ratePrices, nil
}

func (rps *RatePriceService) GetRatePricesByRoomId(ctx *gin.Context, roomID uint64) ([]domain.RatePrice, error) {
	ratePrices, err := rps.repo.GetRatePricesByRoomId(ctx, roomID)
	if err != nil {
        return nil, domain.ErrInternal
    }

    if len(ratePrices) == 0 {
        return nil, domain.ErrDataNotFound
    }

	return ratePrices, nil
}
