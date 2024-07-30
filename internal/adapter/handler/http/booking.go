package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"time"
)

// BookingHandler represents the HTTP handler for booking-related requests
type BookingHandler struct {
	svc port.BookingService
}

// NewBookingHandler creates a new BookingHandler instance
func NewBookingHandler(svc port.BookingService) *BookingHandler {
	return &BookingHandler{
		svc,
	}
}

// createBookingRequest represents the request body for creating a booking
type createBookingRequest struct {
	CustomerID   uint64     `json:"customer_id" binding:"required" example:"1"`
	RoomID       uint64     `json:"room_id" binding:"required" example:"101"`
	CheckInDate  time.Time  `json:"check_in_date" binding:"required" example:"2024-08-01T15:04:05Z"`
	CheckOutDate time.Time  `json:"check_out_date" binding:"required" example:"2024-08-10T15:04:05Z"`
	TotalAmount  float64    `json:"total_amount" binding:"required,gt=0" example:"1000.50"`
}

// CreateBooking godoc
//
//	@Summary		Create a new booking
//	@Description	Create a new booking for a customer
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			createBookingRequest	body		createBookingRequest	true	"Create booking request"
//	@Success		200					{object}	bookingResponse		"Booking created"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/bookings [post]
func (bh *BookingHandler) CreateBooking(ctx *gin.Context) {
	var req createBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	booking := domain.Booking{
		CustomerID:   req.CustomerID,
		RoomID:       req.RoomID,
		CheckInDate:  &req.CheckInDate,
		CheckOutDate: &req.CheckOutDate,
		TotalAmount:  req.TotalAmount,
	}

	createdBooking, err := bh.svc.CreateBooking(ctx, &booking)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newBookingResponse(createdBooking)

	handleSuccess(ctx, rsp)
}

// listBookingsRequest represents the request body for listing bookings
type listBookingsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListBookings godoc
//
//	@Summary		List bookings
//	@Description	List bookings with pagination
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Bookings displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/bookings [get]
//	@Security		BearerAuth
func (bh *BookingHandler) ListBookings(ctx *gin.Context) {
	var req listBookingsRequest
	var bookingsList []bookingResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	bookings, err := bh.svc.ListBookings(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, booking := range bookings {
		bookingsList = append(bookingsList, newBookingResponse(&booking))
	}

	total := uint64(len(bookingsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, bookingsList, "bookings")

	handleSuccess(ctx, rsp)
}

// getBookingRequest represents the request body for getting a booking
type getBookingRequest struct {
	BookingID uint64 `uri:"booking_id" binding:"required,min=1" example:"1"`
}

// GetBooking godoc
//
//	@Summary		Get a booking
//	@Description	Get a booking by id
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Booking ID"
//	@Success		200	{object}	bookingResponse	"Booking displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/bookings/{id} [get]
//	@Security		BearerAuth
func (bh *BookingHandler) GetBooking(ctx *gin.Context) {
	var req getBookingRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	booking, err := bh.svc.GetBooking(ctx, req.BookingID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newBookingResponse(booking)

	handleSuccess(ctx, rsp)
}

// updateBookingRequest represents the request body for updating a booking
type updateBookingRequest struct {
	BookingID    uint64    `json:"booking_id" binding:"required" example:"1"`
	CustomerID   uint64    `json:"customer_id" binding:"required" example:"1"`
	RoomID       uint64    `json:"room_id" binding:"required" example:"101"`
	CheckInDate  time.Time `json:"check_in_date" binding:"required" example:"2024-08-01T15:04:05Z"`
	CheckOutDate time.Time `json:"check_out_date" binding:"required" example:"2024-08-10T15:04:05Z"`
	Status       string    `json:"status" binding:"required" example:"confirmed"`
	TotalAmount  float64   `json:"total_amount" binding:"required,gt=0" example:"1000.50"`
}

// UpdateBooking godoc
//
//	@Summary		Update a booking
//	@Description	Update a booking's details by id
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64				true	"Booking ID"
//	@Param			updateBookingRequest	body		updateBookingRequest	true	"Update booking request"
//	@Success		200					{object}	bookingResponse		"Booking updated"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/bookings/{id} [put]
//	@Security		BearerAuth
func (bh *BookingHandler) UpdateBooking(ctx *gin.Context) {
	var req updateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	booking := domain.Booking{
		ID:           req.BookingID,
		CustomerID:   req.CustomerID,
		RoomID:       req.RoomID,
		CheckInDate:  &req.CheckInDate,
		CheckOutDate: &req.CheckOutDate,
		Status:       req.Status,
		TotalAmount:  req.TotalAmount,
	}

	updatedBooking, err := bh.svc.UpdateBooking(ctx, &booking)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newBookingResponse(updatedBooking)

	handleSuccess(ctx, rsp)
}

// deleteBookingRequest represents the request body for deleting a booking
type deleteBookingRequest struct {
	BookingID uint64 `uri:"booking_id" binding:"required,min=1" example:"1"`
}

// DeleteBooking godoc
//
//	@Summary		Delete a booking
//	@Description	Delete a booking by id
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Booking ID"
//	@Success		200	{object}	response		"Booking deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/bookings/{id} [delete]
//	@Security		BearerAuth
func (bh *BookingHandler) DeleteBooking(ctx *gin.Context) {
	var req deleteBookingRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := bh.svc.DeleteBooking(ctx, req.BookingID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "Booking deleted successfully")
}

// bookingResponse represents the response body for a booking
type bookingResponse struct {
	ID           uint64    `json:"id" example:"1"`
	CustomerID   uint64    `json:"customer_id" example:"1"`
	RoomID       uint64    `json:"room_id" example:"101"`
	CheckInDate  time.Time `json:"check_in_date" example:"2024-08-01T15:04:05Z"`
	CheckOutDate time.Time `json:"check_out_date" example:"2024-08-10T15:04:05Z"`
	Status       string    `json:"status" example:"confirmed"`
	TotalAmount  float64   `json:"total_amount" example:"1000.50"`
	BookingDate  time.Time `json:"booking_date" example:"2024-07-01T15:04:05Z"`
}

// newBookingResponse creates a new booking response
func newBookingResponse(booking *domain.Booking) bookingResponse {
	return bookingResponse{
		ID:           booking.ID,
		CustomerID:   booking.CustomerID,
		RoomID:       booking.RoomID,
		CheckInDate:  *booking.CheckInDate,
		CheckOutDate: *booking.CheckOutDate,
		Status:       booking.Status,
		TotalAmount:  booking.TotalAmount,
		BookingDate:  *booking.BookingDate,
	}
}
