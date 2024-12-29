package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"fmt"
	"time"
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
            "created_bookings",
            "completed_bookings",
            "canceled_bookings",
            "total_amount",
            "status",
        ).
        Values(
            summary.SummaryDate,
            summary.CreatedBookings,
            summary.CompletedBookings,
            summary.CanceledBookings,
            summary.TotalAmount,
            summary.Status,
        ).
        Suffix(`
            ON CONFLICT (summary_date) 
            DO UPDATE SET 
                created_bookings = EXCLUDED.created_bookings,
                completed_bookings = EXCLUDED.completed_bookings,
                canceled_bookings = EXCLUDED.canceled_bookings,
                total_amount = EXCLUDED.total_amount,
                status = EXCLUDED.status,
                updated_at = CURRENT_TIMESTAMP
            RETURNING *
        `)

    sql, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }

    err = dbsr.db.QueryRow(ctx, sql, args...).Scan(
        &summary.SummaryDate,
        &summary.CreatedBookings,
        &summary.CompletedBookings,
        &summary.CanceledBookings,
        &summary.TotalAmount,
        &summary.Status,
        &summary.CreatedAt,
        &summary.UpdatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("error upserting summary: %w", err)
    }

    return summary, nil
}

func (dbsr *DailyBookingSummaryRepository) GetDailyBookingSummaryByDate(ctx *gin.Context, date string) (*domain.DailyBookingSummary, error) {
    var summary domain.DailyBookingSummary

    query := dbsr.db.QueryBuilder.Select(
        "summary_date",
        "created_bookings",
        "completed_bookings",
        "canceled_bookings",
        "total_amount",
        "status",
        "created_at",
        "updated_at",
    ).From("daily_booking_summary").
        Where(sq.Eq{"summary_date": date})

    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("error building query: %w", err)
    }

    err = dbsr.db.QueryRow(ctx, sql, args...).Scan(
        &summary.SummaryDate,
        &summary.CreatedBookings,
        &summary.CompletedBookings,
        &summary.CanceledBookings,
        &summary.TotalAmount,
        &summary.Status,
        &summary.CreatedAt,
        &summary.UpdatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, domain.ErrDataNotFound
        }
        return nil, fmt.Errorf("error querying database: %w", err)
    }

    return &summary, nil
}

func (dbsr *DailyBookingSummaryRepository) ListDailyBookingSummaries(ctx *gin.Context, skip, limit uint64) ([]domain.DailyBookingSummary, uint64, error) {
    var summaries []domain.DailyBookingSummary
    var totalCount uint64

    // Get total count
    countQuery := dbsr.db.QueryBuilder.Select("COUNT(*)").From("daily_booking_summary")
    countSql, countArgs, err := countQuery.ToSql()
    if err != nil {
        return nil, 0, fmt.Errorf("error building count query: %w", err)
    }

    err = dbsr.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
    if err != nil {
        return nil, 0, fmt.Errorf("error counting records: %w", err)
    }

    // Get paginated results
    query := dbsr.db.QueryBuilder.Select(
        "summary_date",
        "created_bookings",
        "completed_bookings",
        "canceled_bookings",
        "total_amount",
        "status",
        "created_at",
        "updated_at",
    ).From("daily_booking_summary").
        OrderBy("summary_date DESC").
        Offset(skip).
        Limit(limit)

    sql, args, err := query.ToSql()
    if err != nil {
        return nil, 0, fmt.Errorf("error building query: %w", err)
    }

    rows, err := dbsr.db.Query(ctx, sql, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("error querying database: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var summary domain.DailyBookingSummary
        err := rows.Scan(
            &summary.SummaryDate,
            &summary.CreatedBookings,
            &summary.CompletedBookings,
            &summary.CanceledBookings,
            &summary.TotalAmount,
            &summary.Status,
            &summary.CreatedAt,
            &summary.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("error scanning row: %w", err)
        }
        summaries = append(summaries, summary)
    }

    return summaries, totalCount, nil
}

func (dbsr *DailyBookingSummaryRepository) UpdateDailyBookingSummary(ctx *gin.Context, summary *domain.DailyBookingSummary) (*domain.DailyBookingSummary, error) {
    query := dbsr.db.QueryBuilder.Update("daily_booking_summary").
        Set("created_bookings", summary.CreatedBookings).
        Set("completed_bookings", summary.CompletedBookings).
        Set("canceled_bookings", summary.CanceledBookings).
        Set("total_amount", summary.TotalAmount).
        Set("status", summary.Status).
        Set("updated_at", time.Now()).
        Where(sq.Eq{"summary_date": summary.SummaryDate}).
        Suffix("RETURNING *")

    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("error building query: %w", err)
    }

    err = dbsr.db.QueryRow(ctx, sql, args...).Scan(
        &summary.SummaryDate,
        &summary.CreatedBookings,
        &summary.CompletedBookings,
        &summary.CanceledBookings,
        &summary.TotalAmount,
        &summary.Status,
        &summary.CreatedAt,
        &summary.UpdatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("error updating record: %w", err)
    }

    return summary, nil
}

func (dbsr *DailyBookingSummaryRepository) DeleteDailyBookingSummary(ctx *gin.Context, date string) error {
    query := dbsr.db.QueryBuilder.Delete("daily_booking_summary").
        Where(sq.Eq{"summary_date": date})

    sql, args, err := query.ToSql()
    if err != nil {
        return fmt.Errorf("error building query: %w", err)
    }

    result, err := dbsr.db.Exec(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("error deleting record: %w", err)
    }

    if result.RowsAffected() == 0 {
        return domain.ErrDataNotFound
    }

    return nil
}
