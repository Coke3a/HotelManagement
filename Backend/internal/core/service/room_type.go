package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type RoomTypeService struct {
	repo    port.RoomTypeRepository
	logRepo port.LogRepository
}

func NewRoomTypeService(repo port.RoomTypeRepository, logRepo port.LogRepository) *RoomTypeService {
	return &RoomTypeService{
		repo,
		logRepo,
	}
}

func (rts *RoomTypeService) CreateRoomType(ctx *gin.Context, roomType *domain.RoomType) (*domain.RoomType, error) {
	if roomType.Name == "" {
		return nil, domain.ErrInvalidData
	}

	roomType, err := rts.repo.CreateRoomType(ctx, roomType)
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
		RecordID:  roomType.ID,
		Action:    "CREATE",
		UserID:    userID.(uint64),
		TableName: "room_types",
	}
	_, err = rts.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return roomType, nil
}

func (rts *RoomTypeService) GetRoomType(ctx *gin.Context, id uint64) (*domain.RoomType, error) {
	roomType, err := rts.repo.GetRoomTypeByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return roomType, nil
}

func (rts *RoomTypeService) ListRoomTypes(ctx *gin.Context, skip, limit uint64) ([]domain.RoomType, uint64, error) {
	roomTypes, totalCount, err := rts.repo.ListRoomTypes(ctx, skip, limit)
	if err != nil {
		return nil, 0, domain.ErrInternal
	}

	return roomTypes, totalCount, nil
}

func (rts *RoomTypeService) UpdateRoomType(ctx *gin.Context, roomType *domain.RoomType) (*domain.RoomType, error) {
	existingRoomType, err := rts.repo.GetRoomTypeByID(ctx, roomType.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	emptyData := roomType.Name == "" &&
		roomType.Description == "" &&
		roomType.Capacity == 0 &&
		roomType.DefaultPrice == 0

	if emptyData {
		return nil, domain.ErrNoUpdatedData
	}

	sameData := existingRoomType.Name == roomType.Name &&
		existingRoomType.Description == roomType.Description &&
		existingRoomType.Capacity == roomType.Capacity &&
		existingRoomType.DefaultPrice == roomType.DefaultPrice

	if sameData {
		return nil, domain.ErrNoUpdatedData
	}

	updatedRoomType, err := rts.repo.UpdateRoomType(ctx, roomType)
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
		RecordID:  roomType.ID,
		Action:    "UPDATE",
		UserID:    userID.(uint64),
		TableName: "room_types",
	}
	_, err = rts.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return updatedRoomType, nil
}

func (rts *RoomTypeService) DeleteRoomType(ctx *gin.Context, id uint64) error {
	_, err := rts.repo.GetRoomTypeByID(ctx, id)
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
		TableName: "room_types",
	}
	_, err = rts.logRepo.CreateLog(ctx, log)
	if err != nil {
		slog.Error("Error creating log", "error", err)
	}

	return rts.repo.DeleteRoomType(ctx, id)
}
