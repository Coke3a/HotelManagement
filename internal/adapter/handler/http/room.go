package http

import (
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
)

// RoomHandler represents the HTTP handler for room-related requests
type RoomHandler struct {
	svc port.RoomService
}

// NewRoomHandler creates a new RoomHandler instance
func NewRoomHandler(svc port.RoomService) *RoomHandler {
	return &RoomHandler{
		svc,
	}
}

// createRoomRequest represents the request body for creating a room
type createRoomRequest struct {
	RoomNumber   string  `json:"room_number" binding:"required" example:"101"`
	Type         string  `json:"type" binding:"required" example:"Deluxe"`
	Description  string  `json:"description" example:"A spacious room with ocean view"`
	Status       string  `json:"status" binding:"required" example:"available"`
	Floor        int     `json:"floor" binding:"required" example:"1"`
	Capacity     int     `json:"capacity" binding:"required,min=1" example:"2"`
	DefaultPrice float64 `json:"default_price" binding:"required,gt=0" example:"150.0"`
}

// CreateRoom godoc
//
//	@Summary		Register a new room
//	@Description	Create a new room in the hotel
//	@Tags			Rooms
//	@Accept			json
//	@Produce		json
//	@Param			createRoomRequest	body		createRoomRequest	true	"Create room request"
//	@Success		200					{object}	roomResponse		"Room registered"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/rooms [post]
func (rh *RoomHandler) CreateRoom(ctx *gin.Context) {
	var req createRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	room := domain.Room{
		RoomNumber:   req.RoomNumber,
		Type:         req.Type,
		Description:  req.Description,
		Status:       req.Status,
		Floor:        req.Floor,
		Capacity:     req.Capacity,
		DefaultPrice: req.DefaultPrice,
	}

	createdRoom, err := rh.svc.RegisterRoom(ctx, &room)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRoomResponse(createdRoom)

	handleSuccess(ctx, rsp)
}

// listRoomsRequest represents the request body for listing rooms
type listRoomsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListRooms godoc
//
//	@Summary		List rooms
//	@Description	List rooms with pagination
//	@Tags			Rooms
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Rooms displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/rooms [get]
//	@Security		BearerAuth
func (rh *RoomHandler) ListRooms(ctx *gin.Context) {
	var req listRoomsRequest
	var roomsList []roomResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	rooms, err := rh.svc.ListRooms(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, room := range rooms {
		roomsList = append(roomsList, newRoomResponse(&room))
	}

	total := uint64(len(roomsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, roomsList, "rooms")

	handleSuccess(ctx, rsp)
}

// getRoomRequest represents the request body for getting a room
type getRoomRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetRoom godoc
//
//	@Summary		Get a room
//	@Description	Get a room by id
//	@Tags			Rooms
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Room ID"
//	@Success		200	{object}	roomResponse	"Room displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/rooms/{id} [get]
//	@Security		BearerAuth
func (rh *RoomHandler) GetRoom(ctx *gin.Context) {
	var req getRoomRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	room, err := rh.svc.GetRoom(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRoomResponse(room)

	handleSuccess(ctx, rsp)
}

// updateRoomRequest represents the request body for updating a room
type updateRoomRequest struct {
	ID           uint64  `json:"id" binding:"required" example:"1"`
	RoomNumber   string  `json:"room_number" binding:"required" example:"101"`
	Type         string  `json:"type" binding:"required" example:"Deluxe"`
	Description  string  `json:"description" example:"A spacious room with ocean view"`
	Status       string  `json:"status" binding:"required" example:"available"`
	Floor        int     `json:"floor" binding:"required" example:"1"`
	Capacity     int     `json:"capacity" binding:"required,min=1" example:"2"`
	DefaultPrice float64 `json:"default_price" binding:"required,gt=0" example:"150.0"`
}

// UpdateRoom godoc
//
//	@Summary		Update a room
//	@Description	Update a room's details by id
//	@Tags			Rooms
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64				true	"Room ID"
//	@Param			updateRoomRequest	body		updateRoomRequest	true	"Update room request"
//	@Success		200					{object}	roomResponse		"Room updated"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/rooms/{id} [put]
//	@Security		BearerAuth
func (rh *RoomHandler) UpdateRoom(ctx *gin.Context) {
	var req updateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	room := domain.Room{
		ID:           req.ID,
		RoomNumber:   req.RoomNumber,
		Type:         req.Type,
		Description:  req.Description,
		Status:       req.Status,
		Floor:        req.Floor,
		Capacity:     req.Capacity,
		DefaultPrice: req.DefaultPrice,
	}

	updatedRoom, err := rh.svc.UpdateRoom(ctx, &room)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRoomResponse(updatedRoom)

	handleSuccess(ctx, rsp)
}

// deleteRoomRequest represents the request body for deleting a room
type deleteRoomRequest struct {
	RoomID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteRoom godoc
//
//	@Summary		Delete a room
//	@Description	Delete a room by id
//	@Tags			Rooms
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Room ID"
//	@Success		200	{object}	response		"Room deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/rooms/{id} [delete]
//	@Security		BearerAuth
func (rh *RoomHandler) DeleteRoom(ctx *gin.Context) {
	var req deleteRoomRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := rh.svc.DeleteRoom(ctx, req.RoomID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "Room deleted successfully")
}

// roomResponse represents the response body for a room
type roomResponse struct {
	ID           uint64    `json:"id" example:"1"`
	RoomNumber   string    `json:"room_number" example:"101"`
	Type         string    `json:"type" example:"Deluxe"`
	Description  string    `json:"description" example:"A spacious room with ocean view"`
	Status       string    `json:"status" example:"available"`
	Floor        int       `json:"floor" example:"1"`
	Capacity     int       `json:"capacity" example:"2"`
	DefaultPrice float64   `json:"default_price" example:"150.0"`
	CreatedAt    time.Time `json:"created_at" example:"2024-07-01T15:04:05Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-07-01T15:04:05Z"`
}

// newRoomResponse creates a new room response
func newRoomResponse(room *domain.Room) roomResponse {
	return roomResponse{
		ID:           room.ID,
		RoomNumber:   room.RoomNumber,
		Type:         room.Type,
		Description:  room.Description,
		Status:       room.Status,
		Floor:        room.Floor,
		Capacity:     room.Capacity,
		DefaultPrice: room.DefaultPrice,
		CreatedAt:    *room.CreatedAt,
		UpdatedAt:    *room.UpdatedAt,
	}
}
