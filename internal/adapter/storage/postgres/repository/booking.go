package repository

import (
	"context"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type BookingRepository struct {
	db *postgres.DB
}

func NewBookingRepository(db *postgres.DB) *BookingRepository {
	return &BookingRepository{
		db,
	}
}

func (br *BookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	query := br.db.QueryBuilder.Insert("bookings").
		Columns("customer_id", "room_id", "check_in_date", "check_out_date", "status", "total_amount", "booking_date").
		Values(booking.CustomerID, booking.RoomID, booking.CheckInDate, booking.CheckOutDate, booking.Status, booking.TotalAmount, booking.BookingDate).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.RoomID,
		&booking.CheckInDate,
		&booking.CheckOutDate,
		&booking.Status,
		&booking.TotalAmount,
		&booking.BookingDate,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errCode := br.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return booking, nil
}

func (br *BookingRepository) GetBookingByID(ctx context.Context, id uint64) (*domain.Booking, error) {
	var booking domain.Booking

	query := br.db.QueryBuilder.Select("*").
		From("bookings").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.RoomID,
		&booking.CheckInDate,
		&booking.CheckOutDate,
		&booking.Status,
		&booking.TotalAmount,
		&booking.BookingDate,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &booking, nil
}

func (br *BookingRepository) ListBookings(ctx context.Context, skip, limit uint64) ([]domain.Booking, error) {
	var bookings []domain.Booking

	query := br.db.QueryBuilder.Select("*").
		From("bookings").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := br.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.CustomerID,
			&booking.RoomID,
			&booking.CheckInDate,
			&booking.CheckOutDate,
			&booking.Status,
			&booking.TotalAmount,
			&booking.BookingDate,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (br *BookingRepository) UpdateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	query := br.db.QueryBuilder.Update("bookings").
		Set("customer_id", sq.Expr("COALESCE(?, customer_id)", booking.CustomerID)).
		Set("room_id", sq.Expr("COALESCE(?, room_id)", booking.RoomID)).
		Set("check_in_date", sq.Expr("COALESCE(?, check_in_date)", booking.CheckInDate)).
		Set("check_out_date", sq.Expr("COALESCE(?, check_out_date)", booking.CheckOutDate)).
		Set("status", sq.Expr("COALESCE(?, status)", booking.Status)).
		Set("total_amount", sq.Expr("COALESCE(?, total_amount)", booking.TotalAmount)).
		Set("booking_date", sq.Expr("COALESCE(?, booking_date)", booking.BookingDate)).
		Where(sq.Eq{"id": booking.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.RoomID,
		&booking.CheckInDate,
		&booking.CheckOutDate,
		&booking.Status,
		&booking.TotalAmount,
		&booking.BookingDate,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errCode := br.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return booking, nil
}

func (br *BookingRepository) DeleteBooking(ctx context.Context, id uint64) error {
	query := br.db.QueryBuilder.Delete("bookings").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = br.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
