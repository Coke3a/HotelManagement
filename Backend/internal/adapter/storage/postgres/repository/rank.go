package repository

import (
	"github.com/gin-gonic/gin"
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type RankRepository struct {
	db *postgres.DB
}

func NewRankRepository(db *postgres.DB) *RankRepository {
	return &RankRepository{
		db,
	}
}

func (rr *RankRepository) CreateRank(ctx *gin.Context, rank *domain.Rank) (*domain.Rank, error) {
	query := rr.db.QueryBuilder.Insert("ranks").
		Columns("rank_name", "description").
		Values(rank.RankName, rank.Description).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rr.db.QueryRow(ctx, sql, args...).Scan(
		&rank.ID,
		&rank.RankName,
		&rank.Description,
	)

	if err != nil {
		if errCode := rr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return rank, nil
}

func (rr *RankRepository) GetRankByID(ctx *gin.Context, id uint64) (*domain.Rank, error) {
	var rank domain.Rank

	query := rr.db.QueryBuilder.Select("*").
		From("ranks").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rr.db.QueryRow(ctx, sql, args...).Scan(
		&rank.ID,
		&rank.RankName,
		&rank.Description,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &rank, nil
}

func (rr *RankRepository) ListRanks(ctx *gin.Context, skip, limit uint64) ([]domain.Rank, error) {
	var ranks []domain.Rank

	query := rr.db.QueryBuilder.Select("*").
		From("ranks").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := rr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rank domain.Rank
		err := rows.Scan(
			&rank.ID,
			&rank.RankName,
			&rank.Description,
		)
		if err != nil {
			return nil, err
		}

		ranks = append(ranks, rank)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ranks, nil
}

func (rr *RankRepository) UpdateRank(ctx *gin.Context, rank *domain.Rank) (*domain.Rank, error) {
	query := rr.db.QueryBuilder.Update("ranks").
		Set("rank_name", sq.Expr("COALESCE(?, rank_name)", rank.RankName)).
		Set("description", sq.Expr("COALESCE(?, description)", rank.Description)).
		Where(sq.Eq{"id": rank.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rr.db.QueryRow(ctx, sql, args...).Scan(
		&rank.ID,
		&rank.RankName,
		&rank.Description,
	)

	if err != nil {
		if errCode := rr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return rank, nil
}

func (rr *RankRepository) DeleteRank(ctx *gin.Context, id uint64) error {
	query := rr.db.QueryBuilder.Delete("ranks").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = rr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
