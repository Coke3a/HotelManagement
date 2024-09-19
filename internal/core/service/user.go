package service

import (
	"context"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/Coke3a/HotelManagement/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo,
	}
}

func (us *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.UserName == "" || user.Email == nil || user.Password == "" {
		return nil, domain.ErrInvalidData
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}
	user.Password = hashedPassword

	// Set the hire date to now
	now := time.Now()
	nullTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	user.LastLogin = &nullTime
	user.HireDate = &now

	createdUser, err := us.repo.CreateUser(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return createdUser, nil
}

func (us *UserService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := us.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	users, err := us.repo.ListUsers(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	isEmpty := user.UserName == "" && user.Email == nil && user.Password == "" && user.Role == nil && user.Rank == nil && user.Status == nil
	isSame := existingUser.UserName == user.UserName && existingUser.Email == user.Email && existingUser.Password == user.Password && existingUser.Role == user.Role && existingUser.Rank == user.Rank && existingUser.Status == user.Status

	if isEmpty || isSame {
		return nil, domain.ErrNoUpdatedData
	}

	if user.Password != "" {
		hashedPassword, err := util.HashPassword(user.Password)
		if err != nil {
			return nil, domain.ErrInternal
		}
		user.Password = hashedPassword
	}

	updatedUser, err := us.repo.UpdateUser(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return updatedUser, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return us.repo.DeleteUser(ctx, id)
}
