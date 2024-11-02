package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type PaymentRepository interface {
	CreatePayment(ctx *gin.Context, payment *domain.Payment) (*domain.Payment, error)
	GetPaymentByID(ctx *gin.Context, id uint64) (*domain.Payment, error)
	ListPayments(ctx *gin.Context, skip, limit uint64) ([]domain.Payment, error)
	UpdatePayment(ctx *gin.Context, payment *domain.Payment) (*domain.Payment, error)
	DeletePayment(ctx *gin.Context, id uint64) error
}

type PaymentService interface {
	ProcessPayment(ctx *gin.Context, payment *domain.Payment) (*domain.Payment, error)
	GetPayment(ctx *gin.Context, id uint64) (*domain.Payment, error)
	ListPayments(ctx *gin.Context, skip, limit uint64) ([]domain.Payment, error)
	UpdatePayment(ctx *gin.Context, payment *domain.Payment) (*domain.Payment, error)
	DeletePayment(ctx *gin.Context, id uint64) error
}
