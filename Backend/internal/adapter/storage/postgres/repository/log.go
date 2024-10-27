package repository

import (
	"context"
	"log/slog"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
)

type LogRepository struct {
	db *postgres.DB
}

func NewLogRepository(db *postgres.DB) *LogRepository {
	return &LogRepository{
		db,
	}
}

func (lr *LogRepository) CreateLog(ctx context.Context, log *domain.Log) (*domain.Log, error) {
	query := lr.db.QueryBuilder.Insert("logs").
		Columns("record_id", "action", "user_id", "table_name").
		Values(log.RecordID, log.Action, log.UserID, log.TableName).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = lr.db.QueryRow(ctx, sql, args...).Scan(
		&log.ID,
		&log.RecordID,
		&log.Action,
		&log.UserID,
		&log.TableName,
		&log.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return log, nil
}
