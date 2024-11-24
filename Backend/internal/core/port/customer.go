package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type CustomerRepository interface {
	CreateCustomer(ctx *gin.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomerByID(ctx *gin.Context, id uint64) (*domain.Customer, error)
	ListCustomers(ctx *gin.Context, skip, limit uint64) ([]domain.Customer, uint64, error)
	UpdateCustomer(ctx *gin.Context, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx *gin.Context, id uint64) error
}

type CustomerService interface {
	RegisterCustomer(ctx *gin.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomer(ctx *gin.Context, id uint64) (*domain.Customer, error)
	ListCustomers(ctx *gin.Context, skip, limit uint64) ([]domain.Customer, uint64, error)
	UpdateCustomer(ctx *gin.Context, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx *gin.Context, id uint64) error
}
