package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CustomerHandler represents the HTTP handler for customer-related requests
type CustomerHandler struct {
	svc port.CustomerService
}

// NewCustomerHandler creates a new CustomerHandler instance
func NewCustomerHandler(svc port.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		svc: svc,
	}
}

// registerCustomerRequest represents the request body for registering a customer
type registerCustomerRequest struct { 
	Name             string    `json:"name" binding:"required"`
	Email            string    `json:"email" binding:"required,email"`
	Phone            string    `json:"phone"`
	Address          string    `json:"address"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Gender           string    `json:"gender"`
	MembershipStatus string    `json:"membership_status"`
	Preferences		 string	`json:"preferences"`
}

// RegisterCustomer handles the HTTP request to register a new customer
func (ch *CustomerHandler) RegisterCustomer(ctx *gin.Context) {
	var req registerCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := domain.Customer{
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		Address:          req.Address,
		DateOfBirth:      req.DateOfBirth,
		Gender:           req.Gender,
		MembershipStatus: req.MembershipStatus,
		Preferences: 	  req.Preferences,
	}

	createdCustomer, err := ch.svc.RegisterCustomer(ctx, &customer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := ch.newCustomerResponse(createdCustomer)

	ctx.JSON(http.StatusOK, resp)
}

// GetCustomer handles the HTTP request to retrieve a customer by ID
func (ch *CustomerHandler) GetCustomer(ctx *gin.Context) {
	customerIDStr, exists := ctx.Params.Get("customer_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "customer_id parameter is required"})
		return
	}

	customerID, err := convertStringToUint64(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer_id format"})
		return
	}

	customer, err := ch.svc.GetCustomer(ctx, customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := ch.newCustomerResponse(customer)

	ctx.JSON(http.StatusOK, resp)
}

// ListCustomers handles the HTTP request to list all customers
func (ch *CustomerHandler) ListCustomers(ctx *gin.Context) {
	var req struct {
		Skip  uint64 `form:"skip,default=0"`
		Limit uint64 `form:"limit,default=10"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customers, err := ch.svc.ListCustomers(ctx, req.Skip, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]customerResponse, len(customers))
	for i, customer := range customers {
		resp[i] = ch.newCustomerResponse(&customer)
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateCustomer handles the HTTP request to update an existing customer
func (ch *CustomerHandler) UpdateCustomer(ctx *gin.Context) {
	var req domain.Customer
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCustomer, err := ch.svc.UpdateCustomer(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := ch.newCustomerResponse(updatedCustomer)

	ctx.JSON(http.StatusOK, resp)
}

// DeleteCustomer handles the HTTP request to delete a customer by ID
func (ch *CustomerHandler) DeleteCustomer(ctx *gin.Context) {
	customerIDStr, exists := ctx.Params.Get("customer_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "customer_id parameter is required"})
		return
	}

	customerID, err := convertStringToUint64(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer_id format"})
		return
	}

	err = ch.svc.DeleteCustomer(ctx, customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// customerResponse represents the response structure for customer-related responses
type customerResponse struct {
	ID               uint64                 `json:"id"`
	Name             string                 `json:"name"`
	Email            string                 `json:"email"`
	Phone            string                 `json:"phone"`
	Address          string                 `json:"address"`
	DateOfBirth      *time.Time              `json:"date_of_birth"`
	Gender           string                 `json:"gender"`
	MembershipStatus string                 `json:"membership_status"`
	JoinDate         *time.Time              `json:"join_date"`
	Preferences      string				 	`json:"preferences"`
	LastVisitDate    *time.Time              `json:"last_visit_date"`
	CreatedAt        *time.Time              `json:"created_at"`
	UpdatedAt        *time.Time              `json:"updated_at"`
}

// newBookingResponse creates a new booking response
func (ch *CustomerHandler) newCustomerResponse(customer *domain.Customer) customerResponse {
	return customerResponse{
		ID:				  customer.ID,
		Name:             customer.Name,
		Email:            customer.Email,
		Phone:            customer.Phone,
		Address:          customer.Address,
		DateOfBirth:      customer.DateOfBirth,
		Gender:           customer.Gender,
		MembershipStatus: customer.MembershipStatus,
		JoinDate:         customer.JoinDate,
		Preferences:      customer.Preferences,
		LastVisitDate:    customer.LastVisitDate,
		CreatedAt:        customer.CreatedAt,
		UpdatedAt:        customer.UpdatedAt,
	}
}
