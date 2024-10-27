package port

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type RankRepository interface {
	CreateRank(ctx context.Context, rank *domain.Rank) (*domain.Rank, error)
	GetRankByID(ctx context.Context, id uint64) (*domain.Rank, error)
	ListRanks(ctx context.Context, skip, limit uint64) ([]domain.Rank, error)
	UpdateRank(ctx context.Context, rank *domain.Rank) (*domain.Rank, error)
	DeleteRank(ctx context.Context, id uint64) error
}

type RankService interface {
	RegisterRank(ctx context.Context, rank *domain.Rank) (*domain.Rank, error)
	GetRank(ctx context.Context, id uint64) (*domain.Rank, error)
	ListRanks(ctx context.Context, skip, limit uint64) ([]domain.Rank, error)
	UpdateRank(ctx context.Context, rank *domain.Rank) (*domain.Rank, error)
	DeleteRank(ctx context.Context, id uint64) error
}
