package service

import (
	"log/slog"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/Coke3a/HotelManagement/internal/core/util"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	repo    port.UserRepository
	logRepo port.LogRepository
}

func NewUserService(repo port.UserRepository, logRepo port.LogRepository) *UserService {
	return &UserService{
		repo,
		logRepo,
	}
}

func (us *UserService) RegisterUser(ctx *gin.Context, user *domain.User) (*domain.User, error) {
	if user.UserName == "" || user.Password == "" {
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

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}

	// Create a log
	log := &domain.Log{
		RecordID:  createdUser.ID,
		Action:    "ADD",
		UserID:    userID.(uint64),
		TableName: "users",
	}

	_, err = us.logRepo.CreateLog(ctx, log)
	if err != nil {
		// Just log the error but don't return it since the user was already created
		slog.Error("Error creating log", "error", err)
	}

	return createdUser, nil
}

func (us *UserService) GetUser(ctx *gin.Context, id uint64) (*domain.User, error) {
	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) ListUsers(ctx *gin.Context, skip, limit uint64) ([]domain.User, uint64, error) {
	users, totalCount, err := us.repo.ListUsers(ctx, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}
	return users, totalCount, nil
}

func (us *UserService) UpdateUser(ctx *gin.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	isEmpty := user.UserName == "" && user.Password == "" && user.Role == 0 && user.Rank == nil && user.Status == nil
	isSame := existingUser.UserName == user.UserName && existingUser.Password == user.Password && existingUser.Role == user.Role && existingUser.Rank == user.Rank && existingUser.Status == user.Status

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

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}

	// Create a log
	log := &domain.Log{
		RecordID:  updatedUser.ID,
		Action:    "UPDATE",
		UserID:    userID.(uint64),
		TableName: "users",
	}
	_, err = us.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return updatedUser, nil
}

func (us *UserService) DeleteUser(ctx *gin.Context, id uint64) error {
	_, err := us.repo.GetUserByID(ctx, id)
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
		TableName: "users",
	}
	_, err = us.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return us.repo.DeleteUser(ctx, id)
}
