package port

import (
    "context"
    "github.com/Coke3a/HotelManagement/internal/core/domain"
)

type RoomTypeRepository interface {
    CreateRoomType(ctx context.Context, roomType *domain.RoomType) (*domain.RoomType, error)
    GetRoomTypeByID(ctx context.Context, id uint64) (*domain.RoomType, error)
    ListRoomTypes(ctx context.Context, skip, limit uint64) ([]domain.RoomType, error)
    UpdateRoomType(ctx context.Context, roomType *domain.RoomType) (*domain.RoomType, error)
    DeleteRoomType(ctx context.Context, id uint64) error
}

type RoomTypeService interface {
    CreateRoomType(ctx context.Context, roomType *domain.RoomType) (*domain.RoomType, error)
    GetRoomType(ctx context.Context, id uint64) (*domain.RoomType, error)
    ListRoomTypes(ctx context.Context, skip, limit uint64) ([]domain.RoomType, error)
    UpdateRoomType(ctx context.Context, roomType *domain.RoomType) (*domain.RoomType, error)
    DeleteRoomType(ctx context.Context, id uint64) error
}