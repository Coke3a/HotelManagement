package port

import (
	"github.com/gin-gonic/gin"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type RankRepository interface {
	CreateRank(ctx *gin.Context, rank *domain.Rank) (*domain.Rank, error)
	GetRankByID(ctx *gin.Context, id uint64) (*domain.Rank, error)
	ListRanks(ctx *gin.Context, skip, limit uint64) ([]domain.Rank, uint64, error)
	UpdateRank(ctx *gin.Context, rank *domain.Rank) (*domain.Rank, error)
	DeleteRank(ctx *gin.Context, id uint64) error
}

type RankService interface {
	RegisterRank(ctx *gin.Context, rank *domain.Rank) (*domain.Rank, error)
	GetRank(ctx *gin.Context, id uint64) (*domain.Rank, error)
	ListRanks(ctx *gin.Context, skip, limit uint64) ([]domain.Rank, uint64, error)
	UpdateRank(ctx *gin.Context, rank *domain.Rank) (*domain.Rank, error)
	DeleteRank(ctx *gin.Context, id uint64) error
}
