package port

import (
	"github.com/gin-gonic/gin"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type RoomTypeRepository interface {
	CreateRoomType(ctx *gin.Context, roomType *domain.RoomType) (*domain.RoomType, error)
	GetRoomTypeByID(ctx *gin.Context, id uint64) (*domain.RoomType, error)
	ListRoomTypes(ctx *gin.Context, skip, limit uint64) ([]domain.RoomType, uint64, error)
	UpdateRoomType(ctx *gin.Context, roomType *domain.RoomType) (*domain.RoomType, error)
	DeleteRoomType(ctx *gin.Context, id uint64) error
}

type RoomTypeService interface {
	CreateRoomType(ctx *gin.Context, roomType *domain.RoomType) (*domain.RoomType, error)
	GetRoomType(ctx *gin.Context, id uint64) (*domain.RoomType, error)
	ListRoomTypes(ctx *gin.Context, skip, limit uint64) ([]domain.RoomType, uint64, error)
	UpdateRoomType(ctx *gin.Context, roomType *domain.RoomType) (*domain.RoomType, error)
	DeleteRoomType(ctx *gin.Context, id uint64) error
}
