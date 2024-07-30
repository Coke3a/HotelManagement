package port

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomerByID(ctx context.Context, id uint64) (*domain.Customer, error)
	ListCustomers(ctx context.Context, skip, limit uint64) ([]domain.Customer, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id uint64) error
}

type CustomerService interface {
	RegisterCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomer(ctx context.Context, id uint64) (*domain.Customer, error)
	ListCustomers(ctx context.Context, skip, limit uint64) ([]domain.Customer, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id uint64) error
}
