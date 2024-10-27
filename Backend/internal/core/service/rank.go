package service

import (
	"context"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
)

type RankService struct {
	repo    port.RankRepository
	logRepo port.LogRepository
}

func NewRankService(repo port.RankRepository, logRepo port.LogRepository) *RankService {
	return &RankService{
		repo,
		logRepo,
	}
}

func (rs *RankService) RegisterRank(ctx context.Context, rank *domain.Rank) (*domain.Rank, error) {
	if rank.RankName == "" {
		return nil, domain.ErrInvalidData
	}

	createdRank, err := rs.repo.CreateRank(ctx, rank)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return createdRank, nil
}

func (rs *RankService) GetRank(ctx context.Context, id uint64) (*domain.Rank, error) {
	rank, err := rs.repo.GetRankByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return rank, nil
}

func (rs *RankService) ListRanks(ctx context.Context, skip, limit uint64) ([]domain.Rank, error) {
	ranks, err := rs.repo.ListRanks(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return ranks, nil
}

func (rs *RankService) UpdateRank(ctx context.Context, rank *domain.Rank) (*domain.Rank, error) {
	existingRank, err := rs.repo.GetRankByID(ctx, rank.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Check if there are changes
	isEmpty := rank.RankName == "" && rank.Description == ""
	isSame := existingRank.RankName == rank.RankName && existingRank.Description == rank.Description

	if isEmpty || isSame {
		return nil, domain.ErrNoUpdatedData
	}

	updatedRank, err := rs.repo.UpdateRank(ctx, rank)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return updatedRank, nil
}

func (rs *RankService) DeleteRank(ctx context.Context, id uint64) error {
	_, err := rs.repo.GetRankByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return rs.repo.DeleteRank(ctx, id)
}
