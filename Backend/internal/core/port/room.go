package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
	"time"
)

type RoomRepository interface {
	CreateRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error)
	GetRoomByID(ctx *gin.Context, id uint64) (*domain.Room, error)
	ListRooms(ctx *gin.Context, skip, limit uint64) ([]domain.Room, uint64, error)
	UpdateRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error)
	DeleteRoom(ctx *gin.Context, id uint64) error
	GetAvailableRooms(ctx *gin.Context, checkInDate, checkOutDate time.Time) ([]domain.RoomWithRoomType, error)
	ListRoomsWithRoomType(ctx *gin.Context, skip, limit uint64) ([]domain.RoomWithRoomType, uint64, error)
}

type RoomService interface {
	RegisterRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error)
	GetRoom(ctx *gin.Context, id uint64) (*domain.Room, error)
	ListRooms(ctx *gin.Context, skip, limit uint64) ([]domain.Room, uint64, error)
	UpdateRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error)
	DeleteRoom(ctx *gin.Context, id uint64) error
	GetAvailableRooms(ctx *gin.Context, checkInDate, checkOutDate time.Time) ([]domain.RoomWithRoomType, error)
	ListRoomsWithRoomType(ctx *gin.Context, skip, limit uint64) ([]domain.RoomWithRoomType, uint64, error)
}
