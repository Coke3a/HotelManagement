package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"errors"
)

// UserHandler handles HTTP requests related to user operations.
// It encapsulates the business logic for user management by utilizing
// the UserService interface from the core layer.
type UserHandler struct {
	svc port.UserService
}

// NewUserHandler creates and returns a new instance of UserHandler.
// It takes a UserService as a parameter to handle the core business logic.
func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

// createUserRequest defines the structure for the request body when creating a new user.
// It includes fields for username, password, email, role, rank, and status.
type createUserRequest struct {
	UserName string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"P@ssw0rd"`
	Role     int    `json:"role" example:"1"`
	Rank     string `json:"rank" example:"Manager"`
	Status   string `json:"status" binding:"required" example:"active"`
}

// CreateUser godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account in the system
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			createUserRequest	body		createUserRequest	true	"User information for registration"
//	@Success		200					{object}	userResponse		"User successfully registered"
//	@Failure		400					{object}	errorResponse		"Invalid input data"
//	@Failure		409					{object}	errorResponse		"User already exists"
//	@Failure		500					{object}	errorResponse		"Server error"
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
		Role:     req.Role,
		Rank:     &req.Rank,
		Status:   &req.Status,
	}

	createdUser, err := uh.svc.RegisterUser(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newUserResponse(createdUser)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, rsp)
}

// ListUsers godoc
//
//	@Summary		Retrieve a list of users
//	@Description	Fetches a paginated list of users from the system
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Number of records to skip"
//	@Param			limit	query		uint64			true	"Maximum number of records to return"
//	@Success		200		{object}	meta			"List of users with pagination metadata"
//	@Failure		400		{object}	errorResponse	"Invalid query parameters"
//	@Failure		500		{object}	errorResponse	"Server error"
//	@Router			/users [get]
//	@Security		BearerAuth
func (uh *UserHandler) ListUsers(ctx *gin.Context) {
    var usersList []userResponse

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

    users, totalCount, err := uh.svc.ListUsers(ctx, skipUint, limitUint)
    if err != nil {
        handleError(ctx, err)
        return
    }

    for _, user := range users {
        rsp, err := newUserResponse(&user)
        if err != nil {
            handleError(ctx, err)
            return
        }
        usersList = append(usersList, rsp)
    }

    meta := newMeta(totalCount, limitUint, skipUint)
    rsp := toMap(meta, usersList, "users")

    handleSuccess(ctx, rsp)
}

// getUserRequest defines the URI parameters for retrieving a specific user.
// It includes a field for the user ID.
type getUserRequest struct {
	UserID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetUser godoc
//
//	@Summary		Retrieve a specific user
//	@Description	Fetches detailed information about a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"User ID"
//	@Success		200	{object}	userResponse	"User details"
//	@Failure		400	{object}	errorResponse	"Invalid user ID"
//	@Failure		404	{object}	errorResponse	"User not found"
//	@Failure		500	{object}	errorResponse	"Server error"
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

	rsp, err := newUserResponse(user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// updateUserRequest defines the structure for the request body when updating a user.
// It includes fields for all updatable user attributes.
type updateUserRequest struct {
	UserID   uint64 `json:"id" binding:"required" example:"1"`
	UserName string `json:"username" example:"johndoe"`
	Role     domain.UserRole `json:"role" example:"1"`
	Rank     string `json:"rank" example:"Manager"`
	Status   string `json:"status" example:"active"`
}

// UpdateUser godoc
//
//	@Summary		Modify an existing user
//	@Description	Updates the details of a specific user identified by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64				true	"User ID"
//	@Param			updateUserRequest	body		updateUserRequest	true	"Updated user information"
//	@Success		200					{object}	userResponse		"Updated user details"
//	@Failure		400					{object}	errorResponse		"Invalid input data"
//	@Failure		404					{object}	errorResponse		"User not found"
//	@Failure		409					{object}	errorResponse		"Update conflict"
//	@Failure		500					{object}	errorResponse		"Server error"
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
		Role:     int(req.Role),
		Rank:     &req.Rank,
		Status:   &req.Status,
	}

	updatedUser, err := uh.svc.UpdateUser(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newUserResponse(updatedUser)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// deleteUserRequest defines the URI parameters for deleting a specific user.
// It includes a field for the user ID.
type deleteUserRequest struct {
	UserID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteUser godoc
//
//	@Summary		Remove a user
//	@Description	Deletes a user from the system by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"User ID"
//	@Success		200	{object}	response		"User successfully deleted"
//	@Failure		400	{object}	errorResponse	"Invalid user ID"
//	@Failure		404	{object}	errorResponse	"User not found"
//	@Failure		500	{object}	errorResponse	"Server error"
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

// userResponse defines the structure for the response body when returning user information.
// It includes all relevant user details.
type userResponse struct {
	ID        uint64    `json:"id" example:"1"`
	UserName  string    `json:"username" example:"johndoe"`
	Role      domain.UserRole    `json:"role" example:"1"`
	Rank      string    `json:"rank" example:"Manager"`
	HireDate  *time.Time `json:"hire_date,omitempty" example:"2024-07-01T15:04:05Z"`
	LastLogin *time.Time `json:"last_login,omitempty" example:"2024-08-01T15:04:05Z"`
	Status    string    `json:"status" example:"active"`
	CreatedAt time.Time `json:"created_at" example:"2024-07-01T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-08-01T15:04:05Z"`
}

// newUserResponse creates and returns a new userResponse from a domain.User object.
// This function is used to convert internal user representations to API responses.
func newUserResponse(user *domain.User) (userResponse, error) {
	if user == nil {
		return userResponse{}, errors.New("user is nil")
	}

	return userResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Role:      domain.UserRole(user.Role),
		Rank:      safeString(user.Rank),
		HireDate:  safeTime(user.HireDate),
		LastLogin: safeTime(user.LastLogin),
		Status:    safeString(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
