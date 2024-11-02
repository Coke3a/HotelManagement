package service

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"log/slog"
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

func (cts *CustomerTypeService) CreateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error) {
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

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  customerType.ID,
		Action:    "CREATE",
		UserID:    userID.(uint64),
		TableName: "customer_types",
	}
	_, err = cts.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return customerType, nil
}

func (cts *CustomerTypeService) GetCustomerType(ctx *gin.Context, id uint64) (*domain.CustomerType, error) {
	customerType, err := cts.repo.GetCustomerTypeByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return customerType, nil
}

func (cts *CustomerTypeService) ListCustomerTypes(ctx *gin.Context, skip, limit uint64) ([]domain.CustomerType, error) {
	customerTypes, err := cts.repo.ListCustomerTypes(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return customerTypes, nil
}

func (cts *CustomerTypeService) UpdateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error) {
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

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  customerType.ID,
		Action:    "UPDATE",
		UserID:    userID.(uint64),
		TableName: "customer_types",
	}
	_, err = cts.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return updatedCustomerType, nil
}

func (cts *CustomerTypeService) DeleteCustomerType(ctx *gin.Context, id uint64) error {
	_, err := cts.repo.GetCustomerTypeByID(ctx, id)
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
		TableName: "customer_types",
	}
	_, err = cts.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return cts.repo.DeleteCustomerType(ctx, id)
}
