package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"errors"
	"strconv"
)

// RatePriceHandler represents the HTTP handler for rate price-related requests
type RatePriceHandler struct {
	svc port.RatePriceService
}

// NewRatePriceHandler creates a new RatePriceHandler instance
func NewRatePriceHandler(svc port.RatePriceService) *RatePriceHandler {
	return &RatePriceHandler{
		svc,
	}
}

// createRatePriceRequest represents the request body for creating a rate price
type createRatePriceRequest struct {
	Name          string  `json:"name" binding:"required" example:"Winter Sale"`
	Description   string  `json:"description" example:"Discount for winter season"`
	PricePerNight float64 `json:"price_per_night" binding:"required,gt=0" example:"10.5"`
	RoomTypeID    uint64  `json:"room_type_id" binding:"required" example:"1"`
}

// CreateRatePrice godoc
//
//	@Summary		Register a new rate price
//	@Description	Create a new rate price for a room
//	@Tags			RatePrices
//	@Accept			json
//	@Produce		json
//	@Param			createRatePriceRequest	body		createRatePriceRequest	true	"Create rate price request"
//	@Success		200						{object}	ratePriceResponse		"Rate price registered"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/rateprices [post]
func (rph *RatePriceHandler) CreateRatePrice(ctx *gin.Context) {
	var req createRatePriceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	ratePrice := domain.RatePrice{
		Name:          req.Name,
		Description:   req.Description,
		PricePerNight: req.PricePerNight,
		RoomTypeID:    req.RoomTypeID,
	}

	createdRatePrice, err := rph.svc.CreateRatePrice(ctx, &ratePrice)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newRatePriceResponse(createdRatePrice)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// listRatePricesRequest represents the request body for listing rate prices
type listRatePricesRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListRatePrices godoc
//
//	@Summary		List rate prices
//	@Description	List rate prices with pagination
//	@Tags			RatePrices
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Rate prices displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/rateprices [get]
//	@Security		BearerAuth
func (rph *RatePriceHandler) ListRatePrices(ctx *gin.Context) {
	var req listRatePricesRequest
	var ratePricesList []ratePriceResponse

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

	ratePrices, totalCount, err := rph.svc.ListRatePrices(ctx, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, ratePrice := range ratePrices {
		rsp, err := newRatePriceResponse(&ratePrice)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ratePricesList = append(ratePricesList, rsp)
	}

	meta := newMeta(totalCount, req.Limit, req.Skip)
	rsp := toMap(meta, ratePricesList, "ratePrices")

	handleSuccess(ctx, rsp)
}

// getRatePriceRequest represents the request body for getting a rate price
type getRatePriceRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetRatePrice godoc
//
//	@Summary		Get a rate price
//	@Description	Get a rate price by id
//	@Tags			RatePrices
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64				true	"Rate Price ID"
//	@Success		200	{object}	ratePriceResponse	"Rate price displayed"
//	@Failure		400	{object}	errorResponse		"Validation error"
//	@Failure		404	{object}	errorResponse		"Data not found error"
//	@Failure		500	{object}	errorResponse		"Internal server error"
//	@Router			/rateprices/{id} [get]
//	@Security		BearerAuth
func (rph *RatePriceHandler) GetRatePrice(ctx *gin.Context) {
	var req getRatePriceRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	ratePrice, err := rph.svc.GetRatePrice(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newRatePriceResponse(ratePrice)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// updateRatePriceRequest represents the request body for updating a rate price
type updateRatePriceRequest struct {
	ID            uint64  `json:"id" binding:"required" example:"1"`
	Name          string  `json:"name" binding:"required" example:"Winter Sale"`
	Description   string  `json:"description" example:"Discount for winter season"`
	PricePerNight float64 `json:"price_per_night" binding:"required" example:"15.5"`
	RoomTypeID    uint64  `json:"room_type_id" binding:"required" example:"1"`
}

// UpdateRatePrice godoc
//
//	@Summary		Update a rate price
//	@Description	Update a rate price's details by id
//	@Tags			RatePrices
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64					true	"Rate Price ID"
//	@Param			updateRatePriceRequest	body		updateRatePriceRequest	true	"Update rate price request"
//	@Success		200					{object}	ratePriceResponse		"Rate price updated"
//	@Failure		400					{object}	errorResponse			"Validation error"
//	@Failure		409					{object}	errorResponse			"Data conflict error"
//	@Failure		500					{object}	errorResponse			"Internal server error"
//	@Router			/rateprices/{id} [put]
//	@Security		BearerAuth
func (rph *RatePriceHandler) UpdateRatePrice(ctx *gin.Context) {
	var req updateRatePriceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	ratePrice := domain.RatePrice{
		ID:            req.ID,
		Name:          req.Name,
		Description:   req.Description,
		PricePerNight: req.PricePerNight,
		RoomTypeID:    req.RoomTypeID,
	}

	updatedRatePrice, err := rph.svc.UpdateRatePrice(ctx, &ratePrice)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newRatePriceResponse(updatedRatePrice)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// deleteRatePriceRequest represents the request body for deleting a rate price
type deleteRatePriceRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteRatePrice godoc
//
//	@Summary		Delete a rate price
//	@Description	Delete a rate price by id
//	@Tags			RatePrices
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Rate Price ID"
//	@Success		200	{object}	response		"Rate price deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/rateprices/{id} [delete]
//	@Security		BearerAuth
func (rph *RatePriceHandler) DeleteRatePrice(ctx *gin.Context) {
	var req deleteRatePriceRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := rph.svc.DeleteRatePrice(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "Rate price deleted successfully")
}

// ratePriceResponse represents the response body for a rate price
type ratePriceResponse struct {
	ID            uint64  `json:"id" example:"1"`
	Name          string  `json:"name" example:"Winter Sale"`
	Description   string  `json:"description" example:"Discount for winter season"`
	PricePerNight float64 `json:"price_per_night" example:"15.5"`
	RoomTypeID    uint64  `json:"room_type_id" example:"101"`
}

// newRatePriceResponse creates a new rate price response
func newRatePriceResponse(ratePrice *domain.RatePrice) (ratePriceResponse, error) {
	if ratePrice == nil {
		return ratePriceResponse{}, errors.New("rate price is nil")
	}

	return ratePriceResponse{
		ID:            ratePrice.ID,
		Name:          ratePrice.Name,
		Description:   ratePrice.Description,
		PricePerNight: ratePrice.PricePerNight,
		RoomTypeID:    ratePrice.RoomTypeID,
	}, nil
}

// GetRatePricesByRoomTypeId godoc
// @Summary Get rate prices by room type ID
// @Description Get a list of rate prices for a specific room type ID
// @Tags rate_prices
// @Accept json
// @Produce json
// @Param room_type_id path uint64 true "Room Type ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rate_prices/by-room-type/{room_type_id} [get]
func (rph *RatePriceHandler) GetRatePricesByRoomTypeId(ctx *gin.Context) {
    roomTypeIDStr := ctx.Param("room_type_id")
    if roomTypeIDStr == "" {
        validationError(ctx, errors.New("room_type_id is required"))
        return
    }

    roomTypeID, err := strconv.ParseUint(roomTypeIDStr, 10, 64)
    if err != nil {
        validationError(ctx, errors.New("invalid room_type_id"))
        return
    }

    ratePrices, totalCount, err := rph.svc.GetRatePricesByRoomTypeId(ctx, roomTypeID)
    if err != nil {
        handleError(ctx, err)
        return
    }

    var ratePricesList []ratePriceResponse
    for _, ratePrice := range ratePrices {
        rsp, err := newRatePriceResponse(&ratePrice)
        if err != nil {
            handleError(ctx, err)
            return
        }
        ratePricesList = append(ratePricesList, rsp)
    }

	meta := newMeta(totalCount, 10, 0)
	rsp := toMap(meta, ratePricesList, "ratePrices")

    handleSuccess(ctx, rsp)
}

// GetRatePricesByRoomId godoc
// @Summary Get rate prices by room ID
// @Description Get a list of rate prices for a specific room ID
// @Tags rate_prices
// @Accept json
// @Produce json
// @Param room_id path uint64 true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rate_prices/by-room/{room_id} [get]
func (rph *RatePriceHandler) GetRatePricesByRoomId(ctx *gin.Context) {
    roomIDStr := ctx.Param("room_id")
    if roomIDStr == "" {
        validationError(ctx, errors.New("room_id is required"))
        return
    }

    roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
    if err != nil {
        validationError(ctx, errors.New("invalid room_id"))
        return
    }

    ratePrices, err := rph.svc.GetRatePricesByRoomId(ctx, roomID)
    if err != nil {
        handleError(ctx, err)
        return
    }

    var ratePricesList []ratePriceResponse
    for _, ratePrice := range ratePrices {
        rsp, err := newRatePriceResponse(&ratePrice)
        if err != nil {
            handleError(ctx, err)
            return
        }
        ratePricesList = append(ratePricesList, rsp)
    }

    handleSuccess(ctx, ratePricesList)
}