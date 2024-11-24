package port

import (
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type LogRepository interface {
	CreateLog(ctx *gin.Context, log *domain.Log) (*domain.Log, error)
	GetLogs(ctx *gin.Context, skip, limit uint64) ([]domain.Log, uint64, error)
}

type LogService interface {
	GetLogs(ctx *gin.Context, skip, limit uint64) ([]domain.Log, uint64, error)
}
