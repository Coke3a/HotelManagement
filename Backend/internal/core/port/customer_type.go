package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type CustomerTypeRepository interface {
    CreateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
    GetCustomerTypeByID(ctx *gin.Context, id uint64) (*domain.CustomerType, error)
    ListCustomerTypes(ctx *gin.Context, skip, limit uint64) ([]domain.CustomerType, error)
    UpdateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
    DeleteCustomerType(ctx *gin.Context, id uint64) error
}

type CustomerTypeService interface {
    CreateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
    GetCustomerType(ctx *gin.Context, id uint64) (*domain.CustomerType, error)
    ListCustomerTypes(ctx *gin.Context, skip, limit uint64) ([]domain.CustomerType, error)
	UpdateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error)
	DeleteCustomerType(ctx *gin.Context, id uint64) error
}
