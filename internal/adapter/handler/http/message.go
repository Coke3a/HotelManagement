package http

import (
	"github.com/Coke3a/TalkPenguin/internal/core/domain"
	"github.com/Coke3a/TalkPenguin/internal/core/port"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	svc port.MessageService
}

func NewMessageHandler(svc port.MessageService) *MessageHandler {
	return &MessageHandler{
		svc,
	}
}

type ExchangingMessageRequest struct {
	UserId       uint64 `json:"user_id" binding:"required" example:"1"`
	ConversId    uint64 `json:"conversation_id" binding:"required" example:"1"`
	MessageText  string `json:"message_text" binding:"required" example:"how are you doing"`
	MessageAudio string `json:"message_audio" example:"{audio file}"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	create a new user account with default role "cashier"
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		registerRequest	true	"Register request"
//	@Success		200				{object}	userResponse	"User created"
//	@Failure		400				{object}	errorResponse	"Validation error"
//	@Failure		401				{object}	errorResponse	"Unauthorized error"
//	@Failure		404				{object}	errorResponse	"Data not found error"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error"
//	@Router			/users [post]
func (ch *MessageHandler) ExchangingMessage(ctx *gin.Context) {
	var req ExchangingMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	message := domain.Message{
		UserId:         req.UserId,
		ConversationId: req.ConversId,
		MessageText:    req.MessageText,
		MessageAudio:   req.MessageAudio,
	}

	_, err := ch.svc.ExchangingMessage(ctx, &message)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp := newExchangingMessageResponse(&message)

	handleSuccess(ctx, rsp)
}

type GetAllMessagesRequest struct {
	UserId    uint64 `uri:"user_id" binding:"required" example:"1"`
	ConversId uint64 `uri:"conversation_id" binding:"required" example:"1"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	create a new user account with default role "cashier"
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		registerRequest	true	"Register request"
//	@Success		200				{object}	userResponse	"User created"
//	@Failure		400				{object}	errorResponse	"Validation error"
//	@Failure		401				{object}	errorResponse	"Unauthorized error"
//	@Failure		404				{object}	errorResponse	"Data not found error"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error"
//	@Router			/users [post]
func (ch *MessageHandler) GetAllMessages(ctx *gin.Context) {
	var req GetAllMessagesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	message := domain.Message{
		UserId:         req.UserId,
		ConversationId: req.ConversId,
	}

	messages, err := ch.svc.GetAllMessagesByConversationId(ctx, message.ConversationId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp := newGetAllMessagesResponse(message.ConversationId, messages)

	handleSuccess(ctx, rsp)
}
