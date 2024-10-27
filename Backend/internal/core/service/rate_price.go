package service

import (
	"context"
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

func (rps *RatePriceService) CreateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error) {
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

	return createdRatePrice, nil
}

func (rps *RatePriceService) GetRatePrice(ctx context.Context, id uint64) (*domain.RatePrice, error) {
	ratePrice, err := rps.repo.GetRatePriceByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return ratePrice, nil
}

func (rps *RatePriceService) ListRatePrices(ctx context.Context, skip, limit uint64) ([]domain.RatePrice, error) {
	ratePrices, err := rps.repo.ListRatePrices(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return ratePrices, nil
}

func (rps *RatePriceService) UpdateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error) {
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

	return updatedRatePrice, nil
}

func (rps *RatePriceService) DeleteRatePrice(ctx context.Context, id uint64) error {
	_, err := rps.repo.GetRatePriceByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return rps.repo.DeleteRatePrice(ctx, id)
}

func (rps *RatePriceService) GetRatePricesByRoomTypeId(ctx context.Context, roomTypeID uint64) ([]domain.RatePrice, error) {
	ratePrices, err := rps.repo.GetRatePricesByRoomTypeId(ctx, roomTypeID)
	if err != nil {
		return nil, domain.ErrInternal
	}

	if len(ratePrices) == 0 {
		return nil, domain.ErrDataNotFound
	}

	return ratePrices, nil
}

func (rps *RatePriceService) GetRatePricesByRoomId(ctx context.Context, roomID uint64) ([]domain.RatePrice, error) {
    ratePrices, err := rps.repo.GetRatePricesByRoomId(ctx, roomID)
    if err != nil {
        return nil, domain.ErrInternal
    }

    if len(ratePrices) == 0 {
        return nil, domain.ErrDataNotFound
    }

    return ratePrices, nil
}