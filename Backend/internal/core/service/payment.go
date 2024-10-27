package service

import (
	"context"
	"time"

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

func (ps *PaymentService) ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	if payment.Amount <= 0 {
		return nil, domain.ErrInvalidData
	}
	if payment.PaymentMethod == domain.PaymentMethodNotSpecified {
		return nil, domain.ErrInvalidData
	}
	if payment.Status == 0 {
		payment.Status = domain.PaymentStatusPending
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

	return createdPayment, nil
}

func (ps *PaymentService) GetPayment(ctx context.Context, id uint64) (*domain.Payment, error) {
	payment, err := ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return payment, nil
}

func (ps *PaymentService) ListPayments(ctx context.Context, skip, limit uint64) ([]domain.Payment, error) {
	payments, err := ps.repo.ListPayments(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return payments, nil
}

func (ps *PaymentService) UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	existingPayment, err := ps.repo.GetPaymentByID(ctx, payment.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	isEmpty := payment.Amount <= 0 && payment.PaymentMethod == domain.PaymentMethodNotSpecified
	isSame := existingPayment.Amount == payment.Amount && existingPayment.PaymentMethod == payment.PaymentMethod && existingPayment.Status == payment.Status

	if isEmpty || isSame {
		return nil, domain.ErrNoUpdatedData
	}

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

	return updatedPayment, nil
}

func (ps *PaymentService) DeletePayment(ctx context.Context, id uint64) error {
	_, err := ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return ps.repo.DeletePayment(ctx, id)
}
