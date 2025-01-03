package http

import (
	"time"
	"errors"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"strconv"
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
	CustomerID   uint64               `json:"customer_id" binding:"required" example:"1"`
	RatePriceId  uint64               `json:"rate_prices_id" binding:"required" example:"1"`
	RoomID       uint64               `json:"room_id" binding:"required" example:"1"`
	RoomTypeID   uint64               `json:"room_type_id" binding:"required" example:"1"`
	CheckInDate  time.Time            `json:"check_in_date" binding:"required" example:"2024-08-01T15:04:05Z"`
	CheckOutDate time.Time            `json:"check_out_date" binding:"required" example:"2024-08-10T15:04:05Z"`
	Status       domain.BookingStatus `json:"status" example:"1"`
	TotalAmount  float64              `json:"total_amount" binding:"required,gt=0" example:"1000.50"`
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

	now := time.Now()
	booking := domain.Booking{
		CustomerID:   req.CustomerID,
		RoomID:       req.RoomID,
		RoomTypeID:   req.RoomTypeID,
		RatePriceId:  req.RatePriceId,
		CheckInDate:  &req.CheckInDate,
		CheckOutDate: &req.CheckOutDate,
		Status:       domain.BookingStatus(req.Status),
		TotalAmount:  req.TotalAmount,
		
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	createdBooking, err := bh.svc.CreateBooking(ctx, &booking)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newBookingResponse(createdBooking)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, rsp)
}

// CreateBookingAndPayment godoc
//
//	@Summary		Create a new booking with payment
//	@Description	Create a new booking for a customer with payment details
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			createBookingRequest	body		createBookingRequest	true	"Create booking request"
//	@Success		200					{object}	bookingResponse		"Booking created with payment"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/bookings/payment [post]
func (bh *BookingHandler) CreateBookingAndPayment(ctx *gin.Context) {
	var req createBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	now := time.Now()
	booking := domain.Booking{
		CustomerID:   req.CustomerID,
		RoomID:       req.RoomID,
		RoomTypeID:   req.RoomTypeID,
		RatePriceId:  req.RatePriceId,
		CheckInDate:  &req.CheckInDate,
		CheckOutDate: &req.CheckOutDate,
		Status:       domain.BookingStatus(req.Status),
		TotalAmount:  req.TotalAmount,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	createdBooking, err := bh.svc.CreateBookingAndPayment(ctx, &booking)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newBookingResponse(createdBooking)
	if err != nil {
		handleError(ctx, err)
		return
	}

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

	bookings, totalCount, err := bh.svc.ListBookings(ctx, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, booking := range bookings {
		rsp, err := newBookingResponse(&booking)
		if err != nil {
			handleError(ctx, err)
			return
		}
		bookingsList = append(bookingsList, rsp)
	}

	meta := newMeta(totalCount, req.Limit, req.Skip)
	rsp := toMap(meta, bookingsList, "bookings")

	handleSuccess(ctx, rsp)
}

// listBookingsRequest represents the request body for listing bookings
type ListBookingsWithFilterRequest struct {
	Skip         uint64                `form:"skip" binding:"required,min=0" example:"0"`
	Limit        uint64                `form:"limit" binding:"required,min=5" example:"5"`
	ID           *uint64               `form:"id,omitempty" example:"1"`
	CustomerID   *uint64               `form:"customer_id,omitempty" example:"2"`
	RatePriceId  *uint64               `form:"rate_price_id,omitempty" example:"3"`
	RoomID       *uint64               `form:"room_id,omitempty" example:"4"`
	RoomTypeID   *uint64               `form:"room_type_id,omitempty" example:"5"`
	CheckInDate  *time.Time            `form:"check_in_date,omitempty" time_format:"2006-01-02" example:"2023-08-31"`
	CheckOutDate *time.Time            `form:"check_out_date,omitempty" time_format:"2006-01-02" example:"2023-09-02"`
	Status       *domain.BookingStatus `form:"status,omitempty" example:"1"`
	TotalAmount  *float64              `form:"total_amount,omitempty" example:"99.99"`
	BookingDate  *time.Time            `form:"booking_date,omitempty" time_format:"2006-01-02" example:"2023-08-01"`
}

// ListBookings godoc
//
//	@Summary		List bookings
//	@Description	List bookings with pagination and filters
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			skip			query		uint64				true	"Skip"
//	@Param			limit			query		uint64				true	"Limit"
//	@Param			id				query		uint64				false	"ID"
//	@Param			customer_id		query		uint64				false	"Customer ID"
//	@Param			rate_price_id	query		uint64				false	"Rate Price ID"
//	@Param			check_in_date	query		string				false	"Check In Date"		time_format:"2006-01-02"
//	@Param			check_out_date	query		string				false	"Check Out Date"	time_format:"2006-01-02"
//	@Param			status			query		uint64				false	"Status"
//	@Param			total_amount	query		float64				false	"Total Amount"
//	@Param			booking_date	query		string				false	"Booking Date"		time_format:"2006-01-02"
//	@Success		200				{object}	meta				"Bookings displayed"
//	@Failure		400				{object}		errorResponse		"Validation error"
//	@Failure		500				{object}		errorResponse		"Internal server error"
//	@Router			/bookings [get]
//	@Security		BearerAuth
func (bh *BookingHandler) ListBookingsWithFilter(ctx *gin.Context) {
	var bookingsList []bookingResponse

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

	// Initialize booking with nil values
	booking := &domain.Booking{}

	// Only parse and set values if they exist in query params
	if id := ctx.Query("id"); id != "" {
		if idUint, err := strconv.ParseUint(id, 10, 64); err == nil {
			booking.ID = idUint
		}
	}

	if customerID := ctx.Query("customer_id"); customerID != "" {
		if customerIDUint, err := strconv.ParseUint(customerID, 10, 64); err == nil {
			booking.CustomerID = customerIDUint
		}
	}

	if checkInDate := ctx.Query("check_in_date"); checkInDate != "" {
		if parsed, err := time.Parse("2006-01-02", checkInDate); err == nil {
			booking.CheckInDate = &parsed
		}
	}

	if checkOutDate := ctx.Query("check_out_date"); checkOutDate != "" {
		if parsed, err := time.Parse("2006-01-02", checkOutDate); err == nil {
			booking.CheckOutDate = &parsed
		}
	}

	if status := ctx.Query("status"); status != "" {
		if statusUint, err := strconv.ParseUint(status, 10, 64); err == nil {
			booking.Status = domain.BookingStatus(statusUint)
		}
	}

	if totalAmount := ctx.Query("total_amount"); totalAmount != "" {
		if amount, err := strconv.ParseFloat(totalAmount, 64); err == nil {
			booking.TotalAmount = amount
		}
	}

	bookings, totalCount, err := bh.svc.ListBookingsWithFilter(ctx, booking, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, booking := range bookings {
		rsp, err := newBookingResponse(&booking)
		if err != nil {
			handleError(ctx, err)
			return
		}
		bookingsList = append(bookingsList, rsp)
	}

	meta := map[string]interface{}{
		"total": totalCount,
		"limit": limitUint,
		"skip":  skipUint,
	}
	
	rsp := map[string]interface{}{
		"bookings": bookingsList,
		"meta":     meta,
	}

	handleSuccess(ctx, rsp)
}

// getBookingRequest represents the request body for getting a booking
type getBookingRequest struct {
	BookingID uint64 `uri:"id" binding:"required,min=1" example:"1"`
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

	rsp, err := newBookingResponse(booking)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// updateBookingRequest represents the request body for updating a booking
type updateBookingRequest struct {
	BookingID    uint64               `json:"id" binding:"required" example:"1"`
	CustomerID   uint64               `json:"customer_id" binding:"required" example:"1"`
	RatePriceId  uint64               `json:"rate_prices_id" binding:"required" example:"1"`
	RoomID       uint64               `json:"room_id" binding:"required" example:"1"`
	RoomTypeID   uint64               `json:"room_type_id" binding:"required" example:"1"`
	CheckInDate  time.Time            `json:"check_in_date" binding:"required" example:"2024-08-01T15:04:05Z"`
	CheckOutDate time.Time            `json:"check_out_date" binding:"required" example:"2024-08-10T15:04:05Z"`
	Status       domain.BookingStatus `json:"status" binding:"required" example:"confirmed"`
	TotalAmount  float64              `json:"total_amount" binding:"required,gt=0" example:"1000.50"`
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

	now := time.Now()
	booking := domain.Booking{
		ID:           req.BookingID,
		CustomerID:   req.CustomerID,
		RatePriceId:  req.RatePriceId,
		RoomID:       req.RoomID,
		RoomTypeID:   req.RoomTypeID,
		CheckInDate:  &req.CheckInDate,
		CheckOutDate: &req.CheckOutDate,
		Status:       domain.BookingStatus(req.Status),
		TotalAmount:  req.TotalAmount,
		UpdatedAt:    &now,
	}

	updatedBooking, err := bh.svc.UpdateBooking(ctx, &booking)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newBookingResponse(updatedBooking)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// deleteBookingRequest represents the request body for deleting a booking
type deleteBookingRequest struct {
	BookingID uint64 `uri:"id" binding:"required,min=1" example:"1"`
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
	ID           uint64               `json:"id" example:"1"`
	CustomerID   uint64               `json:"customer_id" example:"1"`
	RatePriceId  uint64               `json:"rate_prices_id" example:"1"`
	RoomID       uint64               `json:"room_id" example:"1"`
	RoomTypeID   uint64               `json:"room_type_id" example:"1"`
	CheckInDate  time.Time            `json:"check_in_date" example:"2024-08-01T15:04:05Z"`
	CheckOutDate time.Time            `json:"check_out_date" example:"2024-08-10T15:04:05Z"`
	Status       domain.BookingStatus `json:"status" example:"confirmed"`
	TotalAmount  float64              `json:"total_amount" example:"1000.50"`
	CreatedAt    time.Time            `json:"created_at" example:"2024-07-01T15:04:05Z"`
	UpdatedAt    time.Time            `json:"updated_at" example:"2024-07-01T15:04:05Z"`
}

// newBookingResponse creates a new booking response
func newBookingResponse(booking *domain.Booking) (bookingResponse, error) {
	if booking == nil {
		return bookingResponse{}, errors.New("booking is nil")
	}

	var checkInDate, checkOutDate, createdAt, updatedAt time.Time

	if booking.CheckInDate != nil {
		checkInDate = *booking.CheckInDate
	}
	if booking.CheckOutDate != nil {
		checkOutDate = *booking.CheckOutDate
	}
	if booking.CreatedAt != nil {
		createdAt = *booking.CreatedAt
	}
	if booking.UpdatedAt != nil {
		updatedAt = *booking.UpdatedAt
	}

	return bookingResponse{
		ID:           booking.ID,
		CustomerID:   booking.CustomerID,
		RatePriceId:  booking.RatePriceId,
		RoomID:       booking.RoomID,
		RoomTypeID:   booking.RoomTypeID,
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
		Status:       domain.BookingStatus(booking.Status),
		TotalAmount:  booking.TotalAmount,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

type bookingCustomerPaymentResponse struct {
	BookingID         uint64               `json:"booking_id"`
	CustomerID        uint64               `json:"customer_id"`
	BookingPrice      float64              `json:"booking_price"`
	BookingStatus     domain.BookingStatus `json:"booking_status"`
	CheckInDate       string               `json:"check_in_date"`
	CheckOutDate      string               `json:"check_out_date"`
	BookingCreatedAt  string               `json:"booking_created_at"`
	BookingUpdatedAt  string               `json:"booking_updated_at"`
	CustomerFirstName string               `json:"customer_firstname"`
	CustomerSurname   string               `json:"customer_surname"`
	CustomerIdentityNumber string          `json:"customer_identity_number"`
	CustomerAddress string                 `json:"customer_address"`
	RoomID            uint64               `json:"room_id"`
	RoomNumber        string               `json:"room_number"`
	Floor             uint64               `json:"floor"`
	RatePriceID       uint64               `json:"rate_price_id"`
	RoomTypeID        uint64               `json:"room_type_id"`
	RoomTypeName      string               `json:"room_type_name"`
	PaymentID         *uint64              `json:"payment_id"`
	PaymentStatus     *uint64              `json:"payment_status"`
	PaymentUpdateDate *string              `json:"payment_update_date"`
}

func newBookingCustomerPaymentResponse(bcp *domain.BookingCustomerPayment) (*bookingCustomerPaymentResponse, error) {

	if bcp == nil {
		return nil, domain.ErrInvalidData
	}

	response := &bookingCustomerPaymentResponse{
		BookingID:         bcp.BookingID,
		BookingPrice:      bcp.BookingPrice,
		BookingStatus:     bcp.BookingStatus,
		CustomerID:        bcp.CustomerID,
		CustomerFirstName: bcp.CustomerFirstName,
		CustomerSurname:   bcp.CustomerSurname,
		CustomerIdentityNumber: bcp.CustomerIdentityNumber,
		CustomerAddress: bcp.CustomerAddress,
		RoomID:            bcp.RoomID,
		RoomNumber:        bcp.RoomNumber,
		Floor:             bcp.Floor,
		RatePriceID:       bcp.RatePriceID,
		RoomTypeID:        bcp.RoomTypeID,
		RoomTypeName:      bcp.RoomTypeName,
		PaymentStatus:     bcp.PaymentStatus,
	}

	if bcp.CheckInDate != nil {
		response.CheckInDate = bcp.CheckInDate.Format(time.RFC3339)
	}
	if bcp.CheckOutDate != nil {
		response.CheckOutDate = bcp.CheckOutDate.Format(time.RFC3339)
	}
	if bcp.BookingCreatedAt != nil {
		response.BookingCreatedAt = bcp.BookingCreatedAt.Format(time.RFC3339)
	}
	if bcp.BookingUpdatedAt != nil {
		response.BookingUpdatedAt = bcp.BookingUpdatedAt.Format(time.RFC3339)
	}
	if bcp.PaymentID != nil {
		response.PaymentID = bcp.PaymentID
	}
	if bcp.PaymentStatus != nil {
		response.PaymentStatus = bcp.PaymentStatus
	}
	if bcp.PaymentUpdateDate != nil {
		updateDate := bcp.PaymentUpdateDate.Format(time.RFC3339)
		response.PaymentUpdateDate = &updateDate
	}

	return response, nil
}

// GetBookingCustomerPayment godoc
//
//	@Summary		Get a booking customer payment
//	@Description	Get a booking customer payment by booking ID
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Booking ID"
//	@Success		200	{object}	bookingCustomerPaymentResponse	"Booking customer payment displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/bookings/{id}/customer-payment [get]
//	@Security		BearerAuth
func (bh *BookingHandler) GetBookingCustomerPayment(ctx *gin.Context) {
	var req getBookingRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	bcp, err := bh.svc.GetBookingCustomerPayment(ctx, req.BookingID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newBookingCustomerPaymentResponse(bcp)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// ListBookingCustomerPayments godoc
//
//	@Summary		List booking customer payments
//	@Description	List booking customer payments with pagination
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Booking customer payments displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/bookings/customer-payments [get]
//	@Security		BearerAuth
func (bh *BookingHandler) ListBookingCustomerPayments(ctx *gin.Context) {
	var req listBookingsRequest
	var bcpList []bookingCustomerPaymentResponse

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

	bcps, totalCount, err := bh.svc.ListBookingCustomerPayments(ctx, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, bcp := range bcps {
		rsp, err := newBookingCustomerPaymentResponse(&bcp)
		if err != nil {
			handleError(ctx, err)
			return
		}
		bcpList = append(bcpList, *rsp)
	}

	meta := newMeta(totalCount, req.Limit, req.Skip)
	rsp := toMap(meta, bcpList, "booking_customer_payments")

	handleSuccess(ctx, rsp)
}

// listBookingsRequest represents the request body for listing bookings
type ListBookingsWithFilterRequestWithFilter struct {
    Skip              uint64                `form:"skip" binding:"required,min=0" example:"0"`
    Limit             uint64                `form:"limit" binding:"required,min=5" example:"5"`
    BookingID         string                `form:"booking_id,omitempty" example:"1"`
    BookingPrice      string                `form:"booking_price,omitempty" example:"100.00"`
    BookingStatus     string                `form:"booking_status,omitempty" example:"1"`
    CheckInDate       *time.Time            `form:"check_in_date,omitempty" time_format:"2006-01-02" example:"2023-08-01"`
    CheckOutDate      *time.Time            `form:"check_out_date,omitempty" time_format:"2006-01-02" example:"2023-08-01"`
    RoomNumber        string                `form:"room_number,omitempty" example:"101"`
    RoomTypeName      string                `form:"room_type_name,omitempty" example:"Deluxe"`
    CustomerFirstName string                `form:"customer_firstname,omitempty" example:"John"`
    CustomerSurname   string                `form:"customer_surname,omitempty" example:"Doe"`
    PaymentStatus     string                `form:"payment_status,omitempty" example:"1"`
	CreatedAt         *time.Time            `form:"created_at,omitempty" time_format:"2006-01-02" example:"2023-08-01"`
	UpdatedAt         *time.Time            `form:"updated_at,omitempty" time_format:"2006-01-02" example:"2023-08-01"`
}

func (bh *BookingHandler) ListBookingCustomerPaymentsWithFilter(ctx *gin.Context) {
	var bookingsList []bookingCustomerPaymentResponse

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

	// Initialize booking with nil values
	booking := &domain.BookingCustomerPayment{}
	if bookingID := ctx.Query("booking_id"); bookingID != "" {
		if bookingIDUint, err := strconv.ParseUint(bookingID, 10, 64); err == nil {
			booking.BookingID = bookingIDUint
		}
	}
	if bookingPrice := ctx.Query("booking_price"); bookingPrice != "" {
		if bookingPriceFloat, err := strconv.ParseFloat(bookingPrice, 64); err == nil {
			booking.BookingPrice = bookingPriceFloat
		}
	}
	if bookingStatus := ctx.Query("booking_status"); bookingStatus != "" {
		if bookingStatusUint, err := strconv.ParseUint(bookingStatus, 10, 64); err == nil {
			booking.BookingStatus = domain.BookingStatus(bookingStatusUint)
		}
	}
	if checkInDate := ctx.Query("check_in_date"); checkInDate != "" {
		if checkInDate, err := time.Parse("2006-01-02", checkInDate); err == nil {
			booking.CheckInDate = &checkInDate
		}
	}
	if checkOutDate := ctx.Query("check_out_date"); checkOutDate != "" {
		if checkOutDate, err := time.Parse("2006-01-02", checkOutDate); err == nil {
			booking.CheckOutDate = &checkOutDate
		}
	}
	if roomNumber := ctx.Query("room_number"); roomNumber != "" {
		booking.RoomNumber = roomNumber
	}
	if roomTypeName := ctx.Query("room_type_name"); roomTypeName != "" {
		booking.RoomTypeName = roomTypeName
	}
	if customerFirstName := ctx.Query("customer_firstname"); customerFirstName != "" {
		booking.CustomerFirstName = customerFirstName
	}
	if customerSurname := ctx.Query("customer_surname"); customerSurname != "" {
		booking.CustomerSurname = customerSurname
	}
	if paymentStatus := ctx.Query("payment_status"); paymentStatus != "" {
		if paymentStatusUint, err := strconv.ParseUint(paymentStatus, 10, 64); err == nil {
			booking.PaymentStatus = &paymentStatusUint
		}
	}
	if createdAt := ctx.Query("created_at"); createdAt != "" {
		if createdAt, err := time.Parse("2006-01-02", createdAt); err == nil {
			booking.BookingCreatedAt = &createdAt
		}
	}
	if updatedAt := ctx.Query("updated_at"); updatedAt != "" {
		if updatedAt, err := time.Parse("2006-01-02", updatedAt); err == nil {
			booking.BookingUpdatedAt = &updatedAt
		}
	}

	bookings, totalCount, err := bh.svc.ListBookingCustomerPaymentsWithFilter(ctx, booking, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, booking := range bookings {
		rsp, err := newBookingCustomerPaymentResponse(&booking)
		if err != nil {
			handleError(ctx, err)
			return
		}
		bookingsList = append(bookingsList, *rsp)
	}

	meta := map[string]interface{}{
		"total": totalCount,
		"limit": limitUint,
		"skip":  skipUint,
	}
	
	rsp := map[string]interface{}{
		"booking_customer_payments": bookingsList,
		"meta":     meta,
	}

	handleSuccess(ctx, rsp)
}
