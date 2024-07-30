package port

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error)
	GetRoomByID(ctx context.Context, id uint64) (*domain.Room, error)
	ListRooms(ctx context.Context, skip, limit uint64) ([]domain.Room, error)
	UpdateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error)
	DeleteRoom(ctx context.Context, id uint64) error
}

type RoomService interface {
	RegisterRoom(ctx context.Context, room *domain.Room) (*domain.Room, error)
	GetRoom(ctx context.Context, id uint64) (*domain.Room, error)
	ListRooms(ctx context.Context, skip, limit uint64) ([]domain.Room, error)
	UpdateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error)
	DeleteRoom(ctx context.Context, id uint64) error
}
