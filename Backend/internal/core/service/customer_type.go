package service

import (
	"context"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type CustomerTypeService struct {
	repo    port.CustomerTypeRepository
	logRepo port.LogRepository
}

func NewCustomerTypeService(repo port.CustomerTypeRepository, logRepo port.LogRepository) *CustomerTypeService {
	return &CustomerTypeService{
		repo,
		logRepo,
	}
}

func (cts *CustomerTypeService) CreateCustomerType(ctx context.Context, customerType *domain.CustomerType) (*domain.CustomerType, error) {
	if customerType.Name == "" {
		return nil, domain.ErrInvalidData
	}

	customerType, err := cts.repo.CreateCustomerType(ctx, customerType)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	return customerType, nil
}

func (cts *CustomerTypeService) GetCustomerType(ctx context.Context, id uint64) (*domain.CustomerType, error) {
	customerType, err := cts.repo.GetCustomerTypeByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return customerType, nil
}

func (cts *CustomerTypeService) ListCustomerTypes(ctx context.Context, skip, limit uint64) ([]domain.CustomerType, error) {
	customerTypes, err := cts.repo.ListCustomerTypes(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return customerTypes, nil
}

func (cts *CustomerTypeService) UpdateCustomerType(ctx context.Context, customerType *domain.CustomerType) (*domain.CustomerType, error) {
	existingCustomerType, err := cts.repo.GetCustomerTypeByID(ctx, customerType.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	if customerType.Name == "" && customerType.Description == "" {
		return nil, domain.ErrNoUpdatedData
	}

	if existingCustomerType.Name == customerType.Name && existingCustomerType.Description == customerType.Description {
		return nil, domain.ErrNoUpdatedData
	}

	updatedCustomerType, err := cts.repo.UpdateCustomerType(ctx, customerType)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return updatedCustomerType, nil
}

func (cts *CustomerTypeService) DeleteCustomerType(ctx context.Context, id uint64) error {
	_, err := cts.repo.GetCustomerTypeByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return cts.repo.DeleteCustomerType(ctx, id)
}
