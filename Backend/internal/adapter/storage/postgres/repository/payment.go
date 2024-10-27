package repository

import (
	"context"
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type PaymentRepository struct {
	db *postgres.DB
}

func NewPaymentRepository(db *postgres.DB) *PaymentRepository {
	return &PaymentRepository{
		db,
	}
}

func (pr *PaymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	query := pr.db.QueryBuilder.Insert("payments").
		Columns("booking_id", "amount", "payment_method", "payment_date", "status").
		Values(payment.BookingID, payment.Amount, payment.PaymentMethod, payment.PaymentDate, int(payment.Status)).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", sql)

	var status int
	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentDate,
		&status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	payment.Status = domain.PaymentStatus(status)
	return payment, nil
}

func (pr *PaymentRepository) GetPaymentByID(ctx context.Context, id uint64) (*domain.Payment, error) {
	var payment domain.Payment

	query := pr.db.QueryBuilder.Select("*").
		From("payments").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentDate,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &payment, nil
}

func (pr *PaymentRepository) ListPayments(ctx context.Context, skip, limit uint64) ([]domain.Payment, error) {
	var payments []domain.Payment

	query := pr.db.QueryBuilder.Select("*").
		From("payments").
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

	rows, err := pr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.BookingID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentDate,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (pr *PaymentRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	query := pr.db.QueryBuilder.Update("payments").
		Set("amount", sq.Expr("COALESCE(?, amount)", payment.Amount)).
		Set("payment_method", sq.Expr("COALESCE(?, payment_method)", payment.PaymentMethod)).
		Set("payment_date", sq.Expr("COALESCE(?, payment_date)", payment.PaymentDate)).
		Set("status", sq.Expr("COALESCE(?, status)", payment.Status)).
		Where(sq.Eq{"id": payment.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentDate,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errCode := pr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return payment, nil
}

func (pr *PaymentRepository) DeletePayment(ctx context.Context, id uint64) error {
	query := pr.db.QueryBuilder.Delete("payments").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = pr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
