package service

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/Coke3a/HotelManagement/internal/core/util"
	"github.com/gin-gonic/gin"
)

/**
 * AuthService implements port.AuthService interface
 * and provides an access to the user repository
 * and token service
 */
type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

// NewAuthService creates a new auth service instance
func NewAuthService(repo port.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{
		repo,
		ts,
	}
}

// Login gives a registered user an access token and role if the credentials are valid
func (as *AuthService) Login(ctx *gin.Context, userName, password string) (string, domain.UserRole, error) {
	user, err := as.repo.GetUserByUserName(ctx, userName)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return "", 0, domain.ErrInvalidCredentials
		}
		return "", 0, domain.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", 0, domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		return "", 0, domain.ErrTokenCreation
	}

	return accessToken, domain.UserRole(user.Role), nil
}
