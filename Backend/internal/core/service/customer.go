package service

import (
	"context"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type CustomerService struct {
	repo    port.CustomerRepository
	logRepo port.LogRepository
}

func NewCustomerService(repo port.CustomerRepository, logRepo port.LogRepository) *CustomerService {
	return &CustomerService{
		repo,
		logRepo,
	}
}

func (cs *CustomerService) RegisterCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	if customer.FirstName == "" || customer.Surname == "" {
		return nil, domain.ErrInvalidData
	}

	// Add validation for IdentityNumber if needed
	if customer.IdentityNumber == "" {
		return nil, domain.ErrInvalidData
	}

	// Set default values for join date and membership status
	now := time.Now()
	customer.JoinDate = &now

	if customer.CustomerTypeID == 0 {
		customer.CustomerTypeID = 1
	}

	createdCustomer, err := cs.repo.CreateCustomer(ctx, customer)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return createdCustomer, nil
}

func (cs *CustomerService) GetCustomer(ctx context.Context, id uint64) (*domain.Customer, error) {
	customer, err := cs.repo.GetCustomerByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return customer, nil
}

func (cs *CustomerService) ListCustomers(ctx context.Context, skip, limit uint64) ([]domain.Customer, error) {
	customers, err := cs.repo.ListCustomers(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return customers, nil
}

func (cs *CustomerService) UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	existingCustomer, err := cs.repo.GetCustomerByID(ctx, customer.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	isEmpty := customer.FirstName == "" && customer.Surname == "" && customer.Email == "" && customer.Phone == "" && customer.Address == "" && customer.Gender == "" && customer.CustomerTypeID == 0 && customer.Preferences == "" && customer.IdentityNumber == ""
	isSame := existingCustomer.FirstName == customer.FirstName && existingCustomer.Surname == customer.Surname && existingCustomer.Email == customer.Email && existingCustomer.Phone == customer.Phone && existingCustomer.Address == customer.Address && existingCustomer.Gender == customer.Gender && existingCustomer.CustomerTypeID == customer.CustomerTypeID && existingCustomer.Preferences == customer.Preferences && existingCustomer.IdentityNumber == customer.IdentityNumber && existingCustomer.DateOfBirth == customer.DateOfBirth

	if isEmpty || isSame {
		return nil, domain.ErrNoUpdatedData
	}

	updatedCustomer, err := cs.repo.UpdateCustomer(ctx, customer)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return updatedCustomer, nil
}

func (cs *CustomerService) DeleteCustomer(ctx context.Context, id uint64) error {
	_, err := cs.repo.GetCustomerByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return cs.repo.DeleteCustomer(ctx, id)
}