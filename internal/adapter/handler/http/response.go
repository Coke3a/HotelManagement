package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/Coke3a/TalkPenguin/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateConversationResponse struct {
	ConversId         uint64    `json:"convers_id" example:"1"`
	UserId            uint64    `json:"user_id" example:"1"`
	PromptId          uint64    `json:"prompt_id" example:"1"`
	ConversationStart time.Time `json:"start_at" example:"1970-01-01T00:00:00Z"`
}

func newCreateConversationResponse(conversation *domain.Conversation) CreateConversationResponse {
	return CreateConversationResponse{
		ConversId:         conversation.ConversationId,
		UserId:            conversation.UserId,
		PromptId:          conversation.PromptId,
		ConversationStart: *conversation.ConversationStart,
	}
}

type CreateConversationWithMessage struct {
	ConversId         uint64    `json:"convers_id" example:"1"`
	UserId            uint64    `json:"user_id" example:"1"`
	PromptId          uint64    `json:"prompt_id" example:"1"`
	MessageText       string    `json:"message_text" example:"Hello, how are you?"`
	MessageAudio      string    `json:"message_audio" example:"https://example.com/audio.mp3"`
	MessageDate       time.Time `json:"message_date" example:"1970-01-01T00:00:00Z"`
	ConversationStart time.Time `json:"start_at" example:"1970-01-01T00:00:00Z"`
}

func newCreateConversationWithMessageResponse(conversation *domain.Conversation, message *domain.Message) CreateConversationWithMessage {
	return CreateConversationWithMessage{
		ConversId:         conversation.ConversationId,
		UserId:            conversation.UserId,
		PromptId:          conversation.PromptId,
		MessageText:       message.MessageText,
		MessageAudio:      message.MessageAudio,
		MessageDate:       *message.MessageDate,
		ConversationStart: *conversation.ConversationStart,
	}
}

type ExchangingMessageResponse struct {
	MessageText  string    `json:"message_text" example:"Hello, how are you?"`
	MessageAudio string    `json:"message_audio" example:"https://example.com/audio.mp3"`
	MessageDate  time.Time `json:"message_date" example:"1970-01-01T00:00:00Z"`
}

func newExchangingMessageResponse(message *domain.Message) ExchangingMessageResponse {
	var messageDate time.Time
	if message.MessageDate != nil {
		messageDate = *message.MessageDate
	}
	return ExchangingMessageResponse{
		MessageText:  message.MessageText,
		MessageAudio: message.MessageAudio,
		MessageDate:  messageDate,
	}
}

type GetAllMessagesResponse struct {
	ConversId uint64            `json:"convers_id" example:"1"`
	Messages  *[]domain.Message `json:"messages"`
}

func newGetAllMessagesResponse(conversationId uint64, messages *[]domain.Message) GetAllMessagesResponse {
	return GetAllMessagesResponse{
		ConversId: conversationId,
		Messages:  messages,
	}
}

// userResponse represents a user response body
type userResponse struct {
	UserId         uint64    `json:"user_id" example:"1"`
	UserName       string    `json:"name" example:"John Doe"`
	PassWord       string    `json:"password" example:"12345678"`
	Email          string    `json:"email" example:"test@example.com"`
	UserRank       uint64    `json:"user_rank" example:"1"`
	IncorrectLogin uint64    `json:"incorrect_login" example:"0"`
	LastLogin      time.Time `json:"last_login" example:"1970-01-01T00:00:00Z"`
	CreateDate     time.Time `json:"create_date" example:"1970-01-01T00:00:00Z"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		UserId:         user.UserId,
		UserName:       user.UserName,
		PassWord:       user.PassWord,
		Email:          user.Email,
		UserRank:       user.UserRank,
		IncorrectLogin: user.IncorrectLogin,
		LastLogin:      *user.LastLogin,
		CreateDate:     *user.CreateDate,
	}
}

// meta represents metadata for a paginated response
type meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// newMeta is a helper function to create metadata for a paginated response
func newMeta(total, limit, skip uint64) meta {
	return meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	domain.ErrInternal:                   http.StatusInternalServerError,
	domain.ErrDataNotFound:               http.StatusNotFound,
	domain.ErrConflictingData:            http.StatusConflict,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrInvalidToken:               http.StatusUnauthorized,
	domain.ErrExpiredToken:               http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,
	domain.ErrNoUpdatedData:              http.StatusBadRequest,
	domain.ErrInsufficientStock:          http.StatusBadRequest,
	domain.ErrInsufficientPayment:        http.StatusBadRequest,
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// handleAbort sends an error response and aborts the request with the specified status code and error message
func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
