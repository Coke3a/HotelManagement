package service

import (
	"context"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type RatePriceService struct {
	repo port.RatePriceRepository
	// cache port.CacheRepository (if you are using caching)
}

func NewRatePriceService(repo port.RatePriceRepository) *RatePriceService {
	return &RatePriceService{
		repo,
		// cache,
	}
}

func (rps *RatePriceService) RegisterRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error) {
	if ratePrice.Name == "" || ratePrice.DiscountPercentage < 0 || ratePrice.RoomID == 0 {
		return nil, domain.ErrInvalidData
	}

	// Set creation and update timestamps
	now := time.Now()
	ratePrice.CreatedAt = &now
	ratePrice.UpdatedAt = &now

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
		ratePrice.DiscountPercentage == 0 &&
		ratePrice.StartDate == nil &&
		ratePrice.EndDate == nil &&
		ratePrice.RoomID == 0

	sameData := existingRatePrice.Name == ratePrice.Name &&
		existingRatePrice.Description == ratePrice.Description &&
		existingRatePrice.DiscountPercentage == ratePrice.DiscountPercentage &&
		existingRatePrice.StartDate.Equal(*ratePrice.StartDate) &&
		existingRatePrice.EndDate.Equal(*ratePrice.EndDate) &&
		existingRatePrice.RoomID == ratePrice.RoomID

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
