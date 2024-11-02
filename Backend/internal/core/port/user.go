package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	CreateUser(ctx  *gin.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx *gin.Context, id uint64) (*domain.User, error)
	GetUserByUserName(ctx *gin.Context, userName string) (*domain.User, error)
	ListUsers(ctx *gin.Context, skip, limit uint64) ([]domain.User, error)
	UpdateUser(ctx *gin.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx *gin.Context, id uint64) error
}

type UserService interface {
	RegisterUser(ctx  *gin.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx *gin.Context, id uint64) (*domain.User, error)
	ListUsers(ctx *gin.Context, skip, limit uint64) ([]domain.User, error)
	UpdateUser(ctx *gin.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx *gin.Context, id uint64) error
}
