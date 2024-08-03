package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"time"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc port.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

// createUserRequest represents the request body for creating a user
type createUserRequest struct {
	UserName string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"P@ssw0rd"`
	Email    string `json:"email" binding:"required" example:"john.doe@example.com"`
	Role     string `json:"role" binding:"required" example:"admin"`
	Rank     string `json:"rank" example:"Manager"`
	Status   string `json:"status" binding:"required" example:"active"`
}

// CreateUser godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user in the system
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			createUserRequest	body		createUserRequest	true	"Create user request"
//	@Success		200					{object}	userResponse		"User registered"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/users [post]
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	user := domain.User{
		UserName: req.UserName,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
		Rank:     req.Rank,
		Status:   req.Status,
	}

	createdUser, err := uh.svc.RegisterUser(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(createdUser)

	handleSuccess(ctx, rsp)
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListUsers godoc
//
//	@Summary		List users
//	@Description	List users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Users displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/users [get]
//	@Security		BearerAuth
func (uh *UserHandler) ListUsers(ctx *gin.Context) {
	var req listUsersRequest
	var usersList []userResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	users, err := uh.svc.ListUsers(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, user := range users {
		usersList = append(usersList, newUserResponse(&user))
	}

	total := uint64(len(usersList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, usersList, "users")

	handleSuccess(ctx, rsp)
}

// getUserRequest represents the request body for getting a user
type getUserRequest struct {
	UserID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"User ID"
//	@Success		200	{object}	userResponse	"User displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users/{id} [get]
//	@Security		BearerAuth
func (uh *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, err := uh.svc.GetUser(ctx, req.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(user)

	handleSuccess(ctx, rsp)
}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	UserID   uint64 `json:"id" binding:"required" example:"1"`
	UserName string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"NewP@ssw0rd"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Role     string `json:"role" example:"admin"`
	Rank     string `json:"rank" example:"Manager"`
	Status   string `json:"status" example:"active"`
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update a user's details by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64				true	"User ID"
//	@Param			updateUserRequest	body		updateUserRequest	true	"Update user request"
//	@Success		200					{object}	userResponse		"User updated"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/users/{id} [put]
//	@Security		BearerAuth
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := domain.User{
		ID:       req.UserID,
		UserName: req.UserName,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
		Rank:     req.Rank,
		Status:   req.Status,
	}

	updatedUser, err := uh.svc.UpdateUser(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(updatedUser)

	handleSuccess(ctx, rsp)
}

// deleteUserRequest represents the request body for deleting a user
type deleteUserRequest struct {
	UserID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"User ID"
//	@Success		200	{object}	response		"User deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users/{id} [delete]
//	@Security		BearerAuth
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := uh.svc.DeleteUser(ctx, req.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "User deleted successfully")
}

// userResponse represents the response body for a user
type userResponse struct {
	ID        uint64    `json:"id" example:"1"`
	UserName  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Role      string    `json:"role" example:"admin"`
	Rank      string    `json:"rank" example:"Manager"`
	HireDate  time.Time `json:"hire_date" example:"2024-07-01T15:04:05Z"`
	LastLogin time.Time `json:"last_login" example:"2024-08-01T15:04:05Z"`
	Status    string    `json:"status" example:"active"`
	CreatedAt time.Time `json:"created_at" example:"2024-07-01T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-08-01T15:04:05Z"`
}

// newUserResponse creates a new user response
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      user.Role,
		Rank:      user.Rank,
		HireDate:  *user.HireDate,
		LastLogin: *user.LastLogin,
		Status:    user.Status,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}
