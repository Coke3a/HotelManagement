package port

import (
    "context"
    "github.com/Coke3a/HotelManagement/internal/core/domain"
)

type CustomerTypeRepository interface {
    CreateCustomerType(ctx context.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
    GetCustomerTypeByID(ctx context.Context, id uint64) (*domain.CustomerType, error)
    ListCustomerTypes(ctx context.Context, skip, limit uint64) ([]domain.CustomerType, error)
    UpdateCustomerType(ctx context.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
    DeleteCustomerType(ctx context.Context, id uint64) error
}

type CustomerTypeService interface {
    CreateCustomerType(ctx context.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
    GetCustomerType(ctx context.Context, id uint64) (*domain.CustomerType, error)
    ListCustomerTypes(ctx context.Context, skip, limit uint64) ([]domain.CustomerType, error)
    UpdateCustomerType(ctx context.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
    DeleteCustomerType(ctx context.Context, id uint64) error
}
