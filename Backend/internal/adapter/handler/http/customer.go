package http

import (
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"strconv"
	"errors"
)

// CustomerHandler represents the HTTP handler for customer-related requests
type CustomerHandler struct {
	svc port.CustomerService
}

// NewCustomerHandler creates a new CustomerHandler instance
func NewCustomerHandler(svc port.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		svc,
	}
}

// createCustomerRequest represents the request body for creating a customer
type createCustomerRequest struct {
	FirstName       string     `json:"firstname" binding:"required"`
	Surname         string     `json:"surname" binding:"required"`
	IdentityNumber  string     `json:"identity_number" binding:"required"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Address         string     `json:"address"`
	Gender          string     `json:"gender"`
	CustomerTypeID  uint64     `json:"customer_type_id"`
	Preferences     string     `json:"preferences"`
}

// CreateCustomer godoc
//
//	@Summary		Register a new customer
//	@Description	Create a new customer in the system
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			createCustomerRequest	body		createCustomerRequest	true	"Create customer request"
//	@Success		200						{object}	customerResponse		"Customer registered"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/customers [post]
func (ch *CustomerHandler) CreateCustomer(ctx *gin.Context) {
	var req createCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	customer := domain.Customer{
		FirstName:        req.FirstName,
		Surname:          req.Surname,
		IdentityNumber:   req.IdentityNumber,
		Email:            req.Email,
		Phone:            req.Phone,
		Address:          req.Address,
		Gender:           req.Gender,
		CustomerTypeID:   req.CustomerTypeID,
		Preferences:      req.Preferences,
	}

	createdCustomer, err := ch.svc.RegisterCustomer(ctx, &customer)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newCustomerResponse(createdCustomer)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// listCustomersRequest represents the request body for listing customers
type listCustomersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListCustomers godoc
//
//	@Summary		List customers
//	@Description	List customers with pagination
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Customers displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/customers [get]
//	@Security		BearerAuth
func (ch *CustomerHandler) ListCustomers(ctx *gin.Context) {
	var req listCustomersRequest
	var customersList []customerResponse

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

	customers, totalCount, err := ch.svc.ListCustomers(ctx, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, customer := range customers {
		customerResponse, err := newCustomerResponse(&customer)
		if err != nil {
			handleError(ctx, err)
			return
		}
		customersList = append(customersList, customerResponse)
	}

	meta := newMeta(totalCount, req.Limit, req.Skip)
	rsp := toMap(meta, customersList, "customers")

	handleSuccess(ctx, rsp)
}

// getCustomerRequest represents the request body for getting a customer
type getCustomerRequest struct {
	CustomerID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetCustomer godoc
//
//	@Summary		Get a customer
//	@Description	Get a customer by id
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64				true	"Customer ID"
//	@Success		200	{object}	customerResponse	"Customer displayed"
//	@Failure		400	{object}	errorResponse		"Validation error"
//	@Failure		404	{object}	errorResponse		"Data not found error"
//	@Failure		500	{object}	errorResponse		"Internal server error"
//	@Router			/customers/{id} [get]
//	@Security		BearerAuth
func (ch *CustomerHandler) GetCustomer(ctx *gin.Context) {
	var req getCustomerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	customer, err := ch.svc.GetCustomer(ctx, req.CustomerID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newCustomerResponse(customer)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// updateCustomerRequest represents the request body for updating a customer
type updateCustomerRequest struct {
	ID               uint64 `json:"id" binding:"required" example:"1"`
	FirstName        string `json:"firstname" example:"John"`
	Surname          string `json:"surname" example:"Doe"`
	IdentityNumber   string `json:"identity_number" example:"1234567890"`
	Email            string `json:"email" example:"john.doe@example.com"`
	Phone            string `json:"phone" example:"123-456-7890"`
	Address          string `json:"address" example:"123 Elm Street"`
	Gender           string `json:"gender" example:"male"`
	CustomerTypeID   uint64 `json:"customer_type_id" example:"1"`
	Preferences      string `json:"preferences" example:"sea view, non-smoking"`
}

// UpdateCustomer godoc
//
//	@Summary		Update a customer
//	@Description	Update a customer's details by id
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64					true	"Customer ID"
//	@Param			updateCustomerRequest	body		updateCustomerRequest	true	"Update customer request"
//	@Success		200					{object}	customerResponse		"Customer updated"
//	@Failure		400					{object}	errorResponse			"Validation error"
//	@Failure		409					{object}	errorResponse			"Data conflict error"
//	@Failure		500					{object}	errorResponse			"Internal server error"
//	@Router			/customers/{id} [put]
//	@Security		BearerAuth
func (ch *CustomerHandler) UpdateCustomer(ctx *gin.Context) {
	var req updateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	customer := domain.Customer{
		ID:               req.ID,
		FirstName:        req.FirstName,
		Surname:          req.Surname,
		IdentityNumber:   req.IdentityNumber,
		Email:            req.Email,
		Phone:            req.Phone,
		Address:          req.Address,
		Gender:           req.Gender,
		CustomerTypeID:   req.CustomerTypeID,
		Preferences:      req.Preferences,
	}

	updatedCustomer, err := ch.svc.UpdateCustomer(ctx, &customer)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp, err := newCustomerResponse(updatedCustomer)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// deleteCustomerRequest represents the request body for deleting a customer
type deleteCustomerRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteCustomer godoc
//
//	@Summary		Delete a customer
//	@Description	Delete a customer by id
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Customer ID"
//	@Success		200	{object}	response		"Customer deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/customers/{id} [delete]
//	@Security		BearerAuth
func (ch *CustomerHandler) DeleteCustomer(ctx *gin.Context) {
	var req deleteCustomerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := ch.svc.DeleteCustomer(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "Customer deleted successfully")
}

// customerResponse represents the response body for a customer
type customerResponse struct {
	ID              uint64     `json:"id"`
	FirstName       string     `json:"firstname"`
	Surname         string     `json:"surname"`
	IdentityNumber  string     `json:"identity_number"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Address         string     `json:"address"`
	Gender          string     `json:"gender"`
	CustomerTypeID  uint64     `json:"customer_type_id"`
	Preferences     string     `json:"preferences"`
	CreatedAt       string     `json:"created_at,omitempty"`
	UpdatedAt       string     `json:"updated_at,omitempty"`
}

// newCustomerResponse creates a new customer response
func newCustomerResponse(customer *domain.Customer) (customerResponse, error) {
	if customer == nil {
		return customerResponse{}, errors.New("customer is nil")
	}

	var createdAt, updatedAt time.Time

	if customer.CreatedAt != nil {
		createdAt = *customer.CreatedAt
	}
	if customer.UpdatedAt != nil {
		updatedAt = *customer.UpdatedAt
	}

	return customerResponse{
		ID:               customer.ID,
		FirstName:        customer.FirstName,
		Surname:          customer.Surname,
		IdentityNumber:   customer.IdentityNumber,
		Email:            customer.Email,
		Phone:            customer.Phone,
		Address:          customer.Address,
		Gender:           customer.Gender,
		CustomerTypeID:   customer.CustomerTypeID,
		Preferences:      customer.Preferences,
		CreatedAt:        createdAt.Format(time.RFC3339),
		UpdatedAt:        updatedAt.Format(time.RFC3339),
	}, nil
}