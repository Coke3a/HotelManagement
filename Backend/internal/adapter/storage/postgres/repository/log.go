package repository

import (
	"log/slog"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type LogRepository struct {
	db *postgres.DB
}

func NewLogRepository(db *postgres.DB) *LogRepository {
	return &LogRepository{
		db,
	}
}

func (lr *LogRepository) CreateLog(ctx *gin.Context, log *domain.Log) (*domain.Log, error) {
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

func (lr *LogRepository) GetLogs(ctx *gin.Context, skip, limit uint64) ([]domain.Log, error) {
	var logs []domain.Log

	query := lr.db.QueryBuilder.Select("*").
		From("logs").
		OrderBy("id").
		Limit(limit)

	if skip > 0 {
		query = query.Offset(skip)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := lr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log domain.Log
		err := rows.Scan(
			&log.ID,
			&log.TableName,
			&log.RecordID,
			&log.Action,
			&log.UserID,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
