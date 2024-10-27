package service

import (
	"context"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type RoomService struct {
	repo    port.RoomRepository
	logRepo port.LogRepository
}

func NewRoomService(repo port.RoomRepository, logRepo port.LogRepository) *RoomService {
	return &RoomService{
		repo,
		logRepo,
	}
}

func (rs *RoomService) RegisterRoom(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	// Basic validation
	if room.RoomNumber == "" || room.TypeID == 0 {
		return nil, domain.ErrInvalidData
	}

	room, err := rs.repo.CreateRoom(ctx, room)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	return room, nil
}

func (rs *RoomService) GetRoom(ctx context.Context, id uint64) (*domain.Room, error) {
	room, err := rs.repo.GetRoomByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return room, nil
}

func (rs *RoomService) ListRooms(ctx context.Context, skip, limit uint64) ([]domain.Room, error) {
	rooms, err := rs.repo.ListRooms(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return rooms, nil
}

func (rs *RoomService) UpdateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	existingRoom, err := rs.repo.GetRoomByID(ctx, room.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	emptyData := room.RoomNumber == "" &&
		room.TypeID == 0 &&
		room.Description == "" &&
		room.Status == domain.RoomStatus(0) &&
		room.Floor == 0

	sameData := existingRoom.RoomNumber == room.RoomNumber &&
		existingRoom.TypeID == room.TypeID &&
		existingRoom.Description == room.Description &&
		existingRoom.Status == room.Status &&
		existingRoom.Floor == room.Floor

	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	_, err = rs.repo.UpdateRoom(ctx, room)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return room, nil
}

func (rs *RoomService) DeleteRoom(ctx context.Context, id uint64) error {
	_, err := rs.repo.GetRoomByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return rs.repo.DeleteRoom(ctx, id)
}

func (rs *RoomService) GetAvailableRooms(ctx context.Context, checkInDate, checkOutDate time.Time) ([]domain.RoomWithRoomType, error) {
	if checkInDate.After(checkOutDate) {
		return nil, domain.ErrInvalidData
	}

	rooms, err := rs.repo.GetAvailableRooms(ctx, checkInDate, checkOutDate)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return rooms, nil
}

func (rs *RoomService) ListRoomsWithRoomType(ctx context.Context, skip, limit uint64) ([]domain.RoomWithRoomType, error) {
	rooms, err := rs.repo.ListRoomsWithRoomType(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}
	return rooms, nil
}
