package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
)

// RankHandler represents the HTTP handler for rank-related requests
type RankHandler struct {
	svc port.RankService
}

// NewRankHandler creates a new RankHandler instance
func NewRankHandler(svc port.RankService) *RankHandler {
	return &RankHandler{
		svc,
	}
}

// createRankRequest represents the request body for creating a rank
type createRankRequest struct {
	RankName    string `json:"rank_name" binding:"required" example:"Manager"`
	Description string `json:"description" example:"Managerial position"`
}

// CreateRank godoc
//
//	@Summary		Register a new rank
//	@Description	Create a new rank
//	@Tags			Ranks
//	@Accept			json
//	@Produce		json
//	@Param			createRankRequest	body		createRankRequest	true	"Create rank request"
//	@Success		200					{object}	rankResponse		"Rank registered"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/ranks [post]
func (rh *RankHandler) CreateRank(ctx *gin.Context) {
	var req createRankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	rank := domain.Rank{
		RankName:    req.RankName,
		Description: req.Description,
	}

	createdRank, err := rh.svc.RegisterRank(ctx, &rank)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRankResponse(createdRank)

	handleSuccess(ctx, rsp)
}

// listRanksRequest represents the request body for listing ranks
type listRanksRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListRanks godoc
//
//	@Summary		List ranks
//	@Description	List ranks with pagination
//	@Tags			Ranks
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Ranks displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/ranks [get]
//	@Security		BearerAuth
func (rh *RankHandler) ListRanks(ctx *gin.Context) {
	var req listRanksRequest
	var ranksList []rankResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	ranks, totalCount, err := rh.svc.ListRanks(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, rank := range ranks {
		ranksList = append(ranksList, newRankResponse(&rank))
	}

	meta := newMeta(totalCount, req.Limit, req.Skip)
	rsp := toMap(meta, ranksList, "ranks")

	handleSuccess(ctx, rsp)
}

// getRankRequest represents the request body for getting a rank
type getRankRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetRank godoc
//
//	@Summary		Get a rank
//	@Description	Get a rank by id
//	@Tags			Ranks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Rank ID"
//	@Success		200	{object}	rankResponse	"Rank displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/ranks/{id} [get]
//	@Security		BearerAuth
func (rh *RankHandler) GetRank(ctx *gin.Context) {
	var req getRankRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	rank, err := rh.svc.GetRank(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRankResponse(rank)

	handleSuccess(ctx, rsp)
}

// updateRankRequest represents the request body for updating a rank
type updateRankRequest struct {
	ID      uint64 `json:"id" binding:"required" example:"1"`
	RankName    string `json:"rank_name" binding:"required" example:"Manager"`
	Description string `json:"description" example:"Managerial position"`
}

// UpdateRank godoc
//
//	@Summary		Update a rank
//	@Description	Update a rank's details by id
//	@Tags			Ranks
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64				true	"Rank ID"
//	@Param			updateRankRequest	body		updateRankRequest	true	"Update rank request"
//	@Success		200					{object}	rankResponse		"Rank updated"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/ranks/{id} [put]
//	@Security		BearerAuth
func (rh *RankHandler) UpdateRank(ctx *gin.Context) {
	var req updateRankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	rank := domain.Rank{
		ID:          req.ID,
		RankName:    req.RankName,
		Description: req.Description,
	}

	updatedRank, err := rh.svc.UpdateRank(ctx, &rank)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newRankResponse(updatedRank)

	handleSuccess(ctx, rsp)
}

// deleteRankRequest represents the request body for deleting a rank
type deleteRankRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteRank godoc
//
//	@Summary		Delete a rank
//	@Description	Delete a rank by id
//	@Tags			Ranks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Rank ID"
//	@Success		200	{object}	response		"Rank deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/ranks/{id} [delete]
//	@Security		BearerAuth
func (rh *RankHandler) DeleteRank(ctx *gin.Context) {
	var req deleteRankRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := rh.svc.DeleteRank(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "Rank deleted successfully")
}

// rankResponse represents the response body for a rank
type rankResponse struct {
	ID          uint64 `json:"id" example:"1"`
	RankName    string `json:"rank_name" example:"Manager"`
	Description string `json:"description" example:"Managerial position"`
}

// newRankResponse creates a new rank response
func newRankResponse(rank *domain.Rank) rankResponse {
	return rankResponse{
		ID:          rank.ID,
		RankName:    rank.RankName,
		Description: rank.Description,
	}
}
