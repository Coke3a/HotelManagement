package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"strconv"
	"errors"
	"time"
)

type RoomTypeHandler struct {
	svc port.RoomTypeService
}

func NewRoomTypeHandler(svc port.RoomTypeService) *RoomTypeHandler {
	return &RoomTypeHandler{
		svc,
	}
}

type createRoomTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	DefaultPrice float64 `json:"default_price" binding:"required,gt=0"`
}

type updateRoomTypeRequest struct {
	ID          uint64 `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	DefaultPrice float64 `json:"default_price"`
}

type roomTypeResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	DefaultPrice float64 `json:"default_price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (rth *RoomTypeHandler) newRoomTypeResponse(roomType *domain.RoomType) (*roomTypeResponse, error) {
	if roomType == nil {
		return nil, errors.New("roomType is nil")
	}

	if roomType.Description == "" {
		roomType.Description = ""
	}

	return &roomTypeResponse{
		ID:          roomType.ID,
		Name:        roomType.Name,
		Description: roomType.Description,
		Capacity:    roomType.Capacity,
		DefaultPrice: roomType.DefaultPrice,
		CreatedAt:    *roomType.CreatedAt,
		UpdatedAt:    *roomType.UpdatedAt,
	}, nil
}

func (rth *RoomTypeHandler) CreateRoomType(ctx *gin.Context) {
	var req createRoomTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	roomType := &domain.RoomType{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		DefaultPrice: req.DefaultPrice,
	}

	createdRoomType, err := rth.svc.CreateRoomType(ctx, roomType)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := rth.newRoomTypeResponse(createdRoomType)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

func (rth *RoomTypeHandler) GetRoomType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		handleError(ctx, domain.ErrInvalidData)
		return
	}

	roomType, err := rth.svc.GetRoomType(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := rth.newRoomTypeResponse(roomType)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

func (rth *RoomTypeHandler) ListRoomTypes(ctx *gin.Context) {
	skip, _ := strconv.ParseUint(ctx.DefaultQuery("skip", "0"), 10, 64)
	limit, _ := strconv.ParseUint(ctx.DefaultQuery("limit", "10"), 10, 64)

	roomTypes, err := rth.svc.ListRoomTypes(ctx, skip, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	var response []roomTypeResponse
	for _, rt := range roomTypes {
		rsp, err := rth.newRoomTypeResponse(&rt)
		if err != nil {
			handleError(ctx, err)
			return
		}
		response = append(response, *rsp)
	}

	handleSuccess(ctx, response)
}

func (rth *RoomTypeHandler) UpdateRoomType(ctx *gin.Context) {
	var req updateRoomTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, domain.ErrInvalidData)
		return
	}

	roomType := &domain.RoomType{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		DefaultPrice: req.DefaultPrice,
	}

	updatedRoomType, err := rth.svc.UpdateRoomType(ctx, roomType)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := rth.newRoomTypeResponse(updatedRoomType)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

func (rth *RoomTypeHandler) DeleteRoomType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		handleError(ctx, domain.ErrInvalidData)
		return
	}

	err = rth.svc.DeleteRoomType(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Room type deleted successfully"})
}