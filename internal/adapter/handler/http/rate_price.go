package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"time"
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
	Name               string  `json:"name" binding:"required" example:"Winter Sale"`
	Description        string  `json:"description" example:"Discount for winter season"`
	DiscountPercentage float64 `json:"discount_percentage" binding:"required,gt=0" example:"10.5"`
	StartDate          string  `json:"start_date" binding:"required" example:"2024-12-01T00:00:00Z"`
	EndDate            string  `json:"end_date" binding:"required" example:"2025-01-31T23:59:59Z"`
	RoomID             uint64  `json:"room_id" binding:"required" example:"101"`
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

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		validationError(ctx, err)
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		validationError(ctx, err)
		return
	}

	ratePrice := domain.RatePrice{
		Name:              req.Name,
		Description:       req.Description,
		DiscountPercentage: req.DiscountPercentage,
		StartDate:         &startDate,
		EndDate:           &endDate,
		RoomID:            req.RoomID,
	}

	createdRatePrice, err := rph.svc.RegisterRatePrice(ctx, &ratePrice)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRatePriceResponse(createdRatePrice)

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

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	ratePrices, err := rph.svc.ListRatePrices(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, ratePrice := range ratePrices {
		ratePricesList = append(ratePricesList, newRatePriceResponse(&ratePrice))
	}

	total := uint64(len(ratePricesList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, ratePricesList, "ratePrices")

	handleSuccess(ctx, rsp)
}

// getRatePriceRequest represents the request body for getting a rate price
type getRatePriceRequest struct {
	RatePriceID uint64 `uri:"rate_price_id" binding:"required,min=1" example:"1"`
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

	ratePrice, err := rph.svc.GetRatePrice(ctx, req.RatePriceID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRatePriceResponse(ratePrice)

	handleSuccess(ctx, rsp)
}

// updateRatePriceRequest represents the request body for updating a rate price
type updateRatePriceRequest struct {
	RatePriceID       uint64  `json:"rate_price_id" binding:"required" example:"1"`
	Name              string  `json:"name" binding:"required" example:"Winter Sale"`
	Description       string  `json:"description" example:"Discount for winter season"`
	DiscountPercentage float64 `json:"discount_percentage" binding:"required" example:"15.5"`
	StartDate         string  `json:"start_date" binding:"required" example:"2024-12-01T00:00:00Z"`
	EndDate           string  `json:"end_date" binding:"required" example:"2025-01-31T23:59:59Z"`
	RoomID            uint64  `json:"room_id" binding:"required" example:"101"`
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

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		validationError(ctx, err)
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		validationError(ctx, err)
		return
	}

	ratePrice := domain.RatePrice{
		ID:                 req.RatePriceID,
		Name:               req.Name,
		Description:        req.Description,
		DiscountPercentage: req.DiscountPercentage,
		StartDate:          &startDate,
		EndDate:            &endDate,
		RoomID:             req.RoomID,
	}

	updatedRatePrice, err := rph.svc.UpdateRatePrice(ctx, &ratePrice)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRatePriceResponse(updatedRatePrice)

	handleSuccess(ctx, rsp)
}

// deleteRatePriceRequest represents the request body for deleting a rate price
type deleteRatePriceRequest struct {
	RatePriceID uint64 `uri:"rate_price_id" binding:"required,min=1" example:"1"`
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

	err := rph.svc.DeleteRatePrice(ctx, req.RatePriceID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "Rate price deleted successfully")
}

// ratePriceResponse represents the response body for a rate price
type ratePriceResponse struct {
	ID                uint64    `json:"id" example:"1"`
	Name              string    `json:"name" example:"Winter Sale"`
	Description       string    `json:"description" example:"Discount for winter season"`
	DiscountPercentage float64   `json:"discount_percentage" example:"15.5"`
	StartDate         time.Time `json:"start_date" example:"2024-12-01T00:00:00Z"`
	EndDate           time.Time `json:"end_date" example:"2025-01-31T23:59:59Z"`
	RoomID            uint64    `json:"room_id" example:"101"`
}

// newRatePriceResponse creates a new rate price response
func newRatePriceResponse(ratePrice *domain.RatePrice) ratePriceResponse {
	return ratePriceResponse{
		ID:                ratePrice.ID,
		Name:              ratePrice.Name,
		Description:       ratePrice.Description,
		DiscountPercentage: ratePrice.DiscountPercentage,
		StartDate:         *ratePrice.StartDate,
		EndDate:           *ratePrice.EndDate,
		RoomID:            ratePrice.RoomID,
	}
}
