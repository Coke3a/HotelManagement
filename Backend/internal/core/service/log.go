package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type LogService struct {
	repo port.LogRepository
}

func NewLogService(repo port.LogRepository) *LogService {
	return &LogService{
		repo,
	}
}

func (ls *LogService) GetLogs(ctx *gin.Context, skip, limit uint64) ([]domain.Log, error) {
	fmt.Println("skip", skip)
	fmt.Println("limit", limit)
	logs, err := ls.repo.GetLogs(ctx, skip, limit)
	if err != nil {
		fmt.Println("error", err)
		return nil, domain.ErrInternal
	}

	return logs, nil
}
