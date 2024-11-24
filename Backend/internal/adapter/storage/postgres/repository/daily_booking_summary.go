package repository

import (
	"log/slog"
	"github.com/gin-gonic/gin"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type DailyBookingSummaryRepository struct {
	db *postgres.DB
}

func NewDailyBookingSummaryRepository(db *postgres.DB) *DailyBookingSummaryRepository {
	return &DailyBookingSummaryRepository{
		db,
	}
}

func (dbsr *DailyBookingSummaryRepository) CreateDailyBookingSummary(ctx *gin.Context, summary *domain.DailyBookingSummary) (*domain.DailyBookingSummary, error) {
	query := dbsr.db.QueryBuilder.Insert("daily_booking_summary").
		Columns(
			"summary_date",
			"total_bookings",
			"total_amount",
			"pending_bookings",
			"confirmed_bookings",
			"checked_in_bookings",
			"checked_out_bookings",
			"canceled_bookings",
			"completed_bookings",
			"booking_ids",
			"status",
		).
		Values(
			summary.SummaryDate,
			summary.TotalBookings,
			summary.TotalAmount,
			summary.PendingBookings,
			summary.ConfirmedBookings,
			summary.CheckedInBookings,
			summary.CheckedOutBookings,
			summary.CanceledBookings,
			summary.CompletedBookings,
			summary.BookingIDs,
			summary.Status,
		).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = dbsr.db.QueryRow(ctx, sql, args...).Scan(
		&summary.SummaryDate,
		&summary.TotalBookings,
		&summary.TotalAmount,
		&summary.PendingBookings,
		&summary.ConfirmedBookings,
		&summary.CheckedInBookings,
		&summary.CheckedOutBookings,
		&summary.CanceledBookings,
		&summary.CompletedBookings,
		&summary.BookingIDs,
		&summary.Status,
		&summary.CreatedAt,
		&summary.UpdatedAt,
	)

	if err != nil {
		if errCode := dbsr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return summary, nil
}

func (dbsr *DailyBookingSummaryRepository) GetDailyBookingSummaryByDate(ctx *gin.Context, date string) (*domain.DailyBookingSummary, error) {
	var summary domain.DailyBookingSummary

	query := dbsr.db.QueryBuilder.Select("*").
		From("daily_booking_summary").
		Where(sq.Eq{"summary_date": date}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = dbsr.db.QueryRow(ctx, sql, args...).Scan(
		&summary.SummaryDate,
		&summary.TotalBookings,
		&summary.TotalAmount,
		&summary.PendingBookings,
		&summary.ConfirmedBookings,
		&summary.CheckedInBookings,
		&summary.CheckedOutBookings,
		&summary.CanceledBookings,
		&summary.CompletedBookings,
		&summary.BookingIDs,
		&summary.Status,
		&summary.CreatedAt,
		&summary.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &summary, nil
}

func (dbsr *DailyBookingSummaryRepository) ListDailyBookingSummaries(ctx *gin.Context, skip, limit uint64) ([]domain.DailyBookingSummary, uint64, error) {
	var summaries []domain.DailyBookingSummary
	var totalCount uint64

	countQuery := dbsr.db.QueryBuilder.Select("COUNT(*)").From("daily_booking_summary")
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}

	err = dbsr.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := dbsr.db.QueryBuilder.Select("*").
		From("daily_booking_summary").
		OrderBy("summary_date DESC").
		Offset(skip).
		Limit(limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := dbsr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var summary domain.DailyBookingSummary
		err := rows.Scan(
			&summary.SummaryDate,
			&summary.TotalBookings,
			&summary.TotalAmount,
			&summary.PendingBookings,
			&summary.ConfirmedBookings,
			&summary.CheckedInBookings,
			&summary.CheckedOutBookings,
			&summary.CanceledBookings,
			&summary.CompletedBookings,
			&summary.BookingIDs,
			&summary.Status,
			&summary.CreatedAt,
			&summary.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, totalCount, nil
}

func (dbsr *DailyBookingSummaryRepository) UpdateDailyBookingSummary(ctx *gin.Context, summary *domain.DailyBookingSummary) (*domain.DailyBookingSummary, error) {
	query := dbsr.db.QueryBuilder.Update("daily_booking_summary").
		Set("total_bookings", summary.TotalBookings).
		Set("total_amount", summary.TotalAmount).
		Set("pending_bookings", summary.PendingBookings).
		Set("confirmed_bookings", summary.ConfirmedBookings).
		Set("checked_in_bookings", summary.CheckedInBookings).
		Set("checked_out_bookings", summary.CheckedOutBookings).
		Set("canceled_bookings", summary.CanceledBookings).
		Set("completed_bookings", summary.CompletedBookings).
		Set("booking_ids", summary.BookingIDs).
		Set("status", summary.Status).
		Where(sq.Eq{"summary_date": summary.SummaryDate}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = dbsr.db.QueryRow(ctx, sql, args...).Scan(
		&summary.SummaryDate,
		&summary.TotalBookings,
		&summary.TotalAmount,
		&summary.PendingBookings,
		&summary.ConfirmedBookings,
		&summary.CheckedInBookings,
		&summary.CheckedOutBookings,
		&summary.CanceledBookings,
		&summary.CompletedBookings,
		&summary.BookingIDs,
		&summary.Status,
		&summary.CreatedAt,
		&summary.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (dbsr *DailyBookingSummaryRepository) DeleteDailyBookingSummary(ctx *gin.Context, date string) error {
	query := dbsr.db.QueryBuilder.Delete("daily_booking_summary").
		Where(sq.Eq{"summary_date": date})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = dbsr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
