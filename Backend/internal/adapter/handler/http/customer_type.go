package http

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CustomerTypeHandler struct {
	svc port.CustomerTypeService
}

func NewCustomerTypeHandler(svc port.CustomerTypeService) *CustomerTypeHandler {
	return &CustomerTypeHandler{
		svc,
	}
}

type createCustomerTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type updateCustomerTypeRequest struct {
	ID          uint64 `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type customerTypeResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (cth *CustomerTypeHandler) CreateCustomerType(ctx *gin.Context) {
	var req createCustomerTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, domain.ErrInvalidData)
		return
	}

	customerType := &domain.CustomerType{
		Name:        req.Name,
		Description: req.Description,
	}

	createdCustomerType, err := cth.svc.CreateCustomerType(ctx, customerType)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, customerTypeResponse{
		ID:          createdCustomerType.ID,
		Name:        createdCustomerType.Name,
		Description: createdCustomerType.Description,
	})
}

func (cth *CustomerTypeHandler) GetCustomerType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		handleError(ctx, domain.ErrInvalidData)
		return
	}

	customerType, err := cth.svc.GetCustomerType(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, customerTypeResponse{
		ID:          customerType.ID,
		Name:        customerType.Name,
		Description: customerType.Description,
	})
}

func (cth *CustomerTypeHandler) ListCustomerTypes(ctx *gin.Context) {
	skip, _ := strconv.ParseUint(ctx.DefaultQuery("skip", "0"), 10, 64)
	limit, _ := strconv.ParseUint(ctx.DefaultQuery("limit", "10"), 10, 64)

	customerTypes, err := cth.svc.ListCustomerTypes(ctx, skip, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	var response []customerTypeResponse
	for _, ct := range customerTypes {
		response = append(response, customerTypeResponse{
			ID:          ct.ID,
			Name:        ct.Name,
			Description: ct.Description,
		})
	}

	handleSuccess(ctx, response)
}

func (cth *CustomerTypeHandler) UpdateCustomerType(ctx *gin.Context) {
	var req updateCustomerTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, domain.ErrInvalidData)
		return
	}

	customerType := &domain.CustomerType{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	updatedCustomerType, err := cth.svc.UpdateCustomerType(ctx, customerType)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, customerTypeResponse{
		ID:          updatedCustomerType.ID,
		Name:        updatedCustomerType.Name,
		Description: updatedCustomerType.Description,
	})
}

func (cth *CustomerTypeHandler) DeleteCustomerType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		handleError(ctx, domain.ErrInvalidData)
		return
	}

	err = cth.svc.DeleteCustomerType(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Customer type deleted successfully"})
}
