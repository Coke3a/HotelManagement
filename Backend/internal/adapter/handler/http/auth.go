package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

// AuthHandler represents the HTTP handler for authentication-related requests
type AuthHandler struct {
	svc port.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		svc,
	}
}

// loginRequest represents the request body for logging in a user
type loginRequest struct {
	Username    string `json:"username" binding:"required" example:"username1234"`
	Password string `json:"password" binding:"required" example:"12345678" minLength:"1"`
}

// Login godoc
//
//	@Summary		Login and get an access token
//	@Description	Logs in a registered user and returns an access token if the credentials are valid.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request body"
//	@Success		200		{object}	authResponse	"Succesfully logged in"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/users/login [post]
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	token, role, err := ah.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAuthResponse(token, role, req.Username)

	handleSuccess(ctx, rsp)
}

// authResponse represents an authentication response body
type authResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
	Role        domain.UserRole `json:"role" example:"1"`
	Username    string          `json:"username" example:"johndoe"`
}

// newAuthResponse is a helper function to create a response body for handling authentication data
func newAuthResponse(token string, role domain.UserRole, username string) authResponse {
	return authResponse{
		AccessToken: token,
		Role:        role,
		Username:    username,
	}
}