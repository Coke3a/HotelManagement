package http

import (
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"strconv"
	"errors"
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
	TypeID       int     `json:"type_id" binding:"required" example:"1"`
	Description  string  `json:"description" example:"A spacious room with ocean view"`
	Status       int     `json:"status" binding:"required" example:"1"`
	Floor        int     `json:"floor" binding:"required" example:"1"`
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
		TypeID:       req.TypeID,
		Description:  req.Description,
		Status:       domain.RoomStatus(req.Status),
		Floor:        req.Floor,
	}

	createdRoom, err := rh.svc.RegisterRoom(ctx, &room)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newRoomResponse(createdRoom)
	if err != nil {
		handleError(ctx, err)
		return
	}

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

    skip := ctx.Query("skip")
    limit := ctx.Query("limit")

    skipUint, err := strconv.ParseUint(skip, 10, 64)
    if err != nil {
        validationError(ctx, err)
        return
    }

    limitUint, err := strconv.ParseUint(limit, 10, 64)
    if err != nil {
        validationError(ctx, err)
        return
    }

	rooms, err := rh.svc.ListRooms(ctx, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, room := range rooms {
		// roomsList = append(roomsList, newRoomResponse(&room))
		roomResponse, err := newRoomResponse(&room)
		if err != nil {
			handleError(ctx, err)
			return
		}
		roomsList = append(roomsList, roomResponse)
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

	rsp, err := newRoomResponse(room)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// updateRoomRequest represents the request body for updating a room
type updateRoomRequest struct {
	ID           uint64  `json:"id" binding:"required" example:"1"`
	RoomNumber   string  `json:"room_number" binding:"required" example:"101"`
	TypeID       int     `json:"room_type_id" binding:"required" example:"1"`
	Description  string  `json:"description" example:"A spacious room with ocean view"`
	Status       int     `json:"status" binding:"required" example:"1"`
	Floor        int     `json:"floor" binding:"required" example:"1"`
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
		TypeID:       req.TypeID,
		Description:  req.Description,
		Status:       domain.RoomStatus(req.Status),
		Floor:        req.Floor,
	}

	updatedRoom, err := rh.svc.UpdateRoom(ctx, &room)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newRoomResponse(updatedRoom)
	if err != nil {
		handleError(ctx, err)
		return
	}

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
	TypeID       int       `json:"room_type_id" example:"1"`
	Description  string    `json:"description" example:"A spacious room with ocean view"`
	Status       int       `json:"status" example:"1"`
	Floor        int       `json:"floor" example:"1"`
	CreatedAt    time.Time `json:"created_at" example:"2024-07-01T15:04:05Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-07-01T15:04:05Z"`
}

// newRoomResponse creates a new room response
func newRoomResponse(room *domain.Room) (roomResponse, error) {

	if room == nil {
		return roomResponse{}, errors.New("room is nil")
	}

	var createdAt, updatedAt time.Time

	if room.CreatedAt != nil {
		createdAt = *room.CreatedAt
	}
	if room.UpdatedAt != nil {
		updatedAt = *room.UpdatedAt
	}	

	return roomResponse{
		ID:           room.ID,
		RoomNumber:   room.RoomNumber,
		TypeID:       room.TypeID,
		Description:  room.Description,
		Status:       int(room.Status),
		Floor:        room.Floor,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

// roomResponse represents the response body for a room
type availableRoomResponse struct {
	ID           uint64    `json:"id" example:"1"`
	RoomNumber   string    `json:"room_number" example:"101"`
	TypeID       int       `json:"room_type_id" example:"1"`
	TypeName     string    `json:"room_type_name" example:"Deluxe"`
	Description  string    `json:"description" example:"A spacious room with ocean view"`
	Status       int       `json:"status" example:"1"`
	Floor        int       `json:"floor" example:"1"`
	CreatedAt    time.Time `json:"created_at" example:"2024-07-01T15:04:05Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-07-01T15:04:05Z"`
}

func newAvailableRoomResponse(room *domain.RoomWithRoomType) (availableRoomResponse, error) {

	if room == nil {
		return availableRoomResponse{}, errors.New("room is nil")
	}

	var createdAt, updatedAt time.Time

	if room.CreatedAt != nil {
		createdAt = *room.CreatedAt
	}
	if room.UpdatedAt != nil {
		updatedAt = *room.UpdatedAt
	}

	return availableRoomResponse{
		ID:           room.ID,
		RoomNumber:   room.RoomNumber,
		TypeID:       room.TypeID,
		TypeName:     room.TypeName,
		Description:  room.Description,
		Status:       int(room.Status),
		Floor:        room.Floor,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

// GetAvailableRooms godoc
// @Summary Get available rooms
// @Description Get a list of available rooms for a given date range
// @Tags rooms
// @Accept json
// @Produce json
// @Param check_in_date query string true "Check-in date (YYYY-MM-DD)"
// @Param check_out_date query string true "Check-out date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rooms/available [get]
func (rh *RoomHandler) GetAvailableRooms(ctx *gin.Context) {
    checkInDate := ctx.Query("check_in_date")
    checkOutDate := ctx.Query("check_out_date")

    if checkInDate == "" || checkOutDate == "" {
        handleError(ctx, errors.New("check_in_date and check_out_date are required"))
        return
    }

    // Convert string dates to time.Time
    checkIn, err := time.Parse("2006-01-02", checkInDate)
    if err != nil {
        handleError(ctx, errors.New("invalid check-in date format"))
        return
    }

    checkOut, err := time.Parse("2006-01-02", checkOutDate)
    if err != nil {
        handleError(ctx, errors.New("invalid check-out date format"))
        return
    }

    rooms, err := rh.svc.GetAvailableRooms(ctx, checkIn, checkOut)
    if err != nil {
        handleError(ctx, err)
        return
    }

    var roomsList []availableRoomResponse
    for _, room := range rooms {
        roomResponse, err := newAvailableRoomResponse(&room)
        if err != nil {
            handleError(ctx, err)
            return
        }
        roomsList = append(roomsList, roomResponse)
    }

    handleSuccess(ctx, roomsList)
}

func (rh *RoomHandler) ListRoomsWithRoomType(ctx *gin.Context) {

    skip := ctx.Query("skip")
    limit := ctx.Query("limit")

    skipUint, err := strconv.ParseUint(skip, 10, 64)
    if err != nil {
        validationError(ctx, err)
        return
    }

    limitUint, err := strconv.ParseUint(limit, 10, 64)
    if err != nil {
        validationError(ctx, err)
        return
    }
    rooms, err := rh.svc.ListRoomsWithRoomType(ctx, skipUint, limitUint)
    if err != nil {
        handleError(ctx, err)
        return
    }

    var response []availableRoomResponse
    for _, room := range rooms {
        rsp, err := newAvailableRoomResponse(&room)
        if err != nil {
            handleError(ctx, err)
            return
        }
        response = append(response, rsp)
    }

    handleSuccess(ctx, response)
}
