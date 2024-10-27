package http

import (
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"errors"
	"strconv"
)

// PaymentHandler represents the HTTP handler for payment-related requests
type PaymentHandler struct {
	svc port.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler instance
func NewPaymentHandler(svc port.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		svc,
	}
}

// createPaymentRequest represents the request body for creating a payment
type createPaymentRequest struct {
	BookingID     uint64  `json:"booking_id" binding:"required" example:"1"`
	Amount        float64 `json:"amount" binding:"required,gt=0" example:"1000.50"`
	PaymentMethod domain.PaymentMethod  `json:"payment_method" binding:"required" example:"credit_card"`
	Status        domain.PaymentStatus  `json:"status" binding:"required" example:"0"`
}

// CreatePayment godoc
//
//	@Summary		Process a payment
//	@Description	Create a new payment for a booking
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			createPaymentRequest	body		createPaymentRequest	true	"Create payment request"
//	@Success		200					{object}	paymentResponse		"Payment processed"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/payments [post]
func (ph *PaymentHandler) CreatePayment(ctx *gin.Context) {
	var req createPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payment := domain.Payment{
		BookingID:     req.BookingID,
		Amount:        req.Amount,
		PaymentMethod: domain.PaymentMethod(req.PaymentMethod),
		Status:        domain.PaymentStatus(req.Status),
	}

	createdPayment, err := ph.svc.ProcessPayment(ctx, &payment)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newPaymentResponse(createdPayment)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// listPaymentsRequest represents the request body for listing payments
type listPaymentsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListPayments godoc
//
//	@Summary		List payments
//	@Description	List payments with pagination
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Payments displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/payments [get]
//	@Security		BearerAuth
func (ph *PaymentHandler) ListPayments(ctx *gin.Context) {
	var req listPaymentsRequest
	var paymentsList []paymentResponse

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

	payments, err := ph.svc.ListPayments(ctx, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, payment := range payments {
		rsp, err := newPaymentResponse(&payment)
		if err != nil {
			handleError(ctx, err)
			return
		}
		paymentsList = append(paymentsList, rsp)
	}

	total := uint64(len(paymentsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, paymentsList, "payments")

	handleSuccess(ctx, rsp)
}

// getPaymentRequest represents the request body for getting a payment
type getPaymentRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetPayment godoc
//
//	@Summary		Get a payment
//	@Description	Get a payment by id
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Payment ID"
//	@Success		200	{object}	paymentResponse	"Payment displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/payments/{id} [get]
//	@Security		BearerAuth
func (ph *PaymentHandler) GetPayment(ctx *gin.Context) {
	var req getPaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payment, err := ph.svc.GetPayment(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newPaymentResponse(payment)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// updatePaymentRequest represents the request body for updating a payment
type updatePaymentRequest struct {
	ID            uint64  `json:"id" binding:"required" example:"1"`
	Amount        float64 `json:"amount" binding:"required" example:"1000.50"`
	PaymentMethod domain.PaymentMethod  `json:"payment_method" example:"0"`
	Status        int  `json:"status" binding:"required" example:"0"`
}

// UpdatePayment godoc
//
//	@Summary		Update a payment
//	@Description	Update a payment's details by id
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64				true	"Payment ID"
//	@Param			updatePaymentRequest	body		updatePaymentRequest	true	"Update payment request"
//	@Success		200					{object}	paymentResponse		"Payment updated"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/payments/{id} [put]
//	@Security		BearerAuth
func (ph *PaymentHandler) UpdatePayment(ctx *gin.Context) {
	var req updatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payment := domain.Payment{
		ID:            req.ID,
		Amount:        req.Amount,
		PaymentMethod: domain.PaymentMethod(req.PaymentMethod),
		Status:        domain.PaymentStatus(req.Status),
	}

	updatedPayment, err := ph.svc.UpdatePayment(ctx, &payment)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newPaymentResponse(updatedPayment)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// deletePaymentRequest represents the request body for deleting a payment
type deletePaymentRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeletePayment godoc
//
//	@Summary		Delete a payment
//	@Description	Delete a payment by id
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Payment ID"
//	@Success		200	{object}	response		"Payment deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/payments/{id} [delete]
//	@Security		BearerAuth
func (ph *PaymentHandler) DeletePayment(ctx *gin.Context) {
	var req deletePaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := ph.svc.DeletePayment(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "Payment deleted successfully")
}

// paymentResponse represents the response body for a payment
type paymentResponse struct {
	ID            uint64    `json:"id" example:"1"`
	BookingID     uint64    `json:"booking_id" example:"1"`
	Amount        float64   `json:"amount" example:"1000.50"`
	PaymentMethod domain.PaymentMethod    `json:"payment_method" example:"credit_card"`
	PaymentDate   time.Time `json:"payment_date" example:"2024-07-01T15:04:05Z"`
	Status        domain.PaymentStatus    `json:"status" example:"0"`
}

// newPaymentResponse creates a new payment response
func newPaymentResponse(payment *domain.Payment) (paymentResponse, error) {
	if payment == nil {
		return paymentResponse{}, errors.New("payment is nil")
	}

	var paymentDate time.Time
	if payment.PaymentDate != nil {
		paymentDate = *payment.PaymentDate
	}

	return paymentResponse{
		ID:            payment.ID,
		BookingID:     payment.BookingID,
		Amount:        payment.Amount,
		PaymentMethod: domain.PaymentMethod(payment.PaymentMethod),
		PaymentDate:   paymentDate,
		Status:        payment.Status,
	}, nil
}
