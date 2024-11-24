package service

import (
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"log/slog"
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

func (rs *RoomService) RegisterRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error) {
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

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  room.ID,
		Action:    "CREATE",
		UserID:    userID.(uint64),
		TableName: "rooms",
	}
	_, err = rs.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return room, nil
}

func (rs *RoomService) GetRoom(ctx *gin.Context, id uint64) (*domain.Room, error) {
	room, err := rs.repo.GetRoomByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return room, nil
}

func (rs *RoomService) ListRooms(ctx *gin.Context, skip, limit uint64) ([]domain.Room, uint64, error) {
	rooms, totalCount, err := rs.repo.ListRooms(ctx, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}

	return rooms, totalCount, nil
}

func (rs *RoomService) UpdateRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error) {
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

	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, domain.ErrUnauthorized
	}
	// Create a log
	log := &domain.Log{
		RecordID:  room.ID,
		Action:    "UPDATE",
		UserID:    userID.(uint64),
		TableName: "rooms",
	}
	_, err = rs.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return room, nil
}

func (rs *RoomService) DeleteRoom(ctx *gin.Context, id uint64) error {
	_, err := rs.repo.GetRoomByID(ctx, id)
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
		TableName: "rooms",
	}
	_, err = rs.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return rs.repo.DeleteRoom(ctx, id)
}

func (rs *RoomService) GetAvailableRooms(ctx *gin.Context, checkInDate, checkOutDate time.Time) ([]domain.RoomWithRoomType, error) {
	if checkInDate.After(checkOutDate) {
		return nil, domain.ErrInvalidData
	}

	rooms, err := rs.repo.GetAvailableRooms(ctx, checkInDate, checkOutDate)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return rooms, nil
}

func (rs *RoomService) ListRoomsWithRoomType(ctx *gin.Context, skip, limit uint64) ([]domain.RoomWithRoomType, uint64, error) {
	rooms, totalCount, err := rs.repo.ListRoomsWithRoomType(ctx, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}
	return rooms, totalCount, nil
}
