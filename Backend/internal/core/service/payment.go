package service

import (
	"github.com/gin-gonic/gin"
	"time"
	"log/slog"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type PaymentService struct {
	repo    port.PaymentRepository
	logRepo port.LogRepository
}

func NewPaymentService(repo port.PaymentRepository, logRepo port.LogRepository) *PaymentService {
	return &PaymentService{
		repo,
		logRepo,
	}
}

func (ps *PaymentService) ProcessPayment(ctx *gin.Context, payment *domain.Payment) (*domain.Payment, error) {
	if payment.Amount <= 0 {
		return nil, domain.ErrInvalidData
	}
	if payment.PaymentMethod == domain.PaymentMethodNotSpecified {
		return nil, domain.ErrInvalidData
	}
	if payment.Status == 0 {
		payment.Status = domain.PaymentStatusUnpaid
	}

	// Set the payment date if it's not already set
	if payment.PaymentDate == nil {
		now := time.Now()
		payment.PaymentDate = &now
	}

	createdPayment, err := ps.repo.CreatePayment(ctx, payment)
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
		RecordID:  createdPayment.ID,
		Action:    "CREATE",
		UserID:    userID.(uint64),
		TableName: "payments",
	}
	_, err = ps.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return createdPayment, nil
}

func (ps *PaymentService) GetPayment(ctx *gin.Context, id uint64) (*domain.Payment, error) {
	payment, err := ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return payment, nil
}

func (ps *PaymentService) ListPayments(ctx *gin.Context, skip, limit uint64) ([]domain.Payment, uint64, error) {
	payments, totalCount, err := ps.repo.ListPayments(ctx, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}

	return payments, totalCount, nil
}

func (ps *PaymentService) UpdatePayment(ctx *gin.Context, payment *domain.Payment) (*domain.Payment, error) {
	// existingPayment, err := ps.repo.GetPaymentByID(ctx, payment.ID)
	// if err != nil {
	// 	if err == domain.ErrDataNotFound {
	// 		return nil, err
	// 	}
	// 	return nil, domain.ErrInternal
	// }

	// // Check if there are changes
	// isEmpty := payment.Amount <= 0 && payment.PaymentMethod == domain.PaymentMethodNotSpecified
	// isSame := existingPayment.Amount == payment.Amount && existingPayment.PaymentMethod == payment.PaymentMethod && existingPayment.Status == payment.Status

	// if isEmpty || isSame {
	// 	return nil, domain.ErrNoUpdatedData
	// }

	// Update timestamp
	now := time.Now()
	payment.PaymentDate = &now
	
	updatedPayment, err := ps.repo.UpdatePayment(ctx, payment)
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
		RecordID:  payment.ID,
		Action:    "UPDATE",
		UserID:    userID.(uint64),
		TableName: "payments",
	}
	_, err = ps.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return updatedPayment, nil
}

func (ps *PaymentService) DeletePayment(ctx *gin.Context, id uint64) error {
	_, err := ps.repo.GetPaymentByID(ctx, id)
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
		TableName: "payments",
	}
	_, err = ps.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return ps.repo.DeletePayment(ctx, id)
}
