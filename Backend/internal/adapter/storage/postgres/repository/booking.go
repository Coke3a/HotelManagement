package repository

import (
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/gin-gonic/gin"
	"fmt"
)

type BookingRepository struct {
	db *postgres.DB
}

func NewBookingRepository(db *postgres.DB) *BookingRepository {
	return &BookingRepository{
		db,
	}
}

func (br *BookingRepository) CreateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error) {
	query := br.db.QueryBuilder.Insert("bookings").
		Columns("customer_id", "rate_prices_id", "room_id", "room_type_id", "check_in_date", "check_out_date", "status", "total_amount", "booking_date").
		Values(
			booking.CustomerID,
			booking.RatePriceId,
			booking.RoomID,
			booking.RoomTypeID,
			booking.CheckInDate.Format("2006-01-02"),
			booking.CheckOutDate.Format("2006-01-02"),
			booking.Status,
			booking.TotalAmount,
			booking.BookingDate,
		).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.RatePriceId,
		&booking.RoomID,
		&booking.RoomTypeID,
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

func (br *BookingRepository) GetBookingByID(ctx *gin.Context, id uint64) (*domain.Booking, error) {
	var booking domain.Booking

	query := br.db.QueryBuilder.Select("*").
		From("bookings").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.RatePriceId,
		&booking.RoomID,
		&booking.RoomTypeID,
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

func (br *BookingRepository) ListBookings(ctx *gin.Context, skip, limit uint64) ([]domain.Booking, uint64, error) {
	var bookings []domain.Booking
	var totalCount uint64

	countQuery := br.db.QueryBuilder.Select("COUNT(*)").From("bookings")
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = br.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := br.db.QueryBuilder.Select("*").
		OrderBy("id").
		Limit(limit)

	if skip > 0 {
		query = query.Offset(skip)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := br.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.CustomerID,
			&booking.RatePriceId,
			&booking.RoomID,
			&booking.RoomTypeID,
			&booking.CheckInDate,
			&booking.CheckOutDate,
			&booking.Status,
			&booking.TotalAmount,
			&booking.BookingDate,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return bookings, totalCount, nil
}

func (rr *BookingRepository) ListBookingsWithFilter(ctx *gin.Context, booking *domain.Booking, skip, limit uint64) ([]domain.Booking, uint64, error) {
	var bookings []domain.Booking
	var totalCount uint64

	countQuery := rr.db.QueryBuilder.Select("COUNT(*)").From("bookings")
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = rr.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := rr.db.QueryBuilder.Select("*").
		From("bookings").
		OrderBy("id").
		Limit(limit)

	if skip > 0 {
		query = query.Offset(skip)
	}

	// Add WHERE clauses for each non-nil field in the Booking struct
	if booking.ID != 0 {
		query = query.Where("id = ?", booking.ID)
	}
	if booking.CustomerID != 0 {
		query = query.Where("customer_id = ?", booking.CustomerID)
	}
	if booking.RatePriceId != 0 {
		query = query.Where("rate_price_id = ?", booking.RatePriceId)
	}
	if booking.RoomID != 0 {
		query = query.Where("room_id = ?", booking.RoomID)
	}
	if booking.RoomTypeID != 0 {
		query = query.Where("room_type_id = ?", booking.RoomTypeID)
	}
	if booking.CheckInDate != nil {
		query = query.Where("check_in_date = ?", *booking.CheckInDate)
	}
	if booking.CheckOutDate != nil {
		query = query.Where("check_out_date = ?", *booking.CheckOutDate)
	}
	if booking.Status != 0 {
		query = query.Where("status = ?", booking.Status)
	}
	if booking.TotalAmount != 0 {
		query = query.Where("total_amount = ?", booking.TotalAmount)
	}
	if booking.BookingDate != nil {
		query = query.Where("booking_date = ?", *booking.BookingDate)
	}
	if booking.CreatedAt != nil {
		query = query.Where("created_at = ?", *booking.CreatedAt)
	}
	if booking.UpdatedAt != nil {
		query = query.Where("updated_at = ?", *booking.UpdatedAt)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := rr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.CustomerID,
			&booking.RatePriceId,
			&booking.RoomID,
			&booking.RoomTypeID,
			&booking.CheckInDate,
			&booking.CheckOutDate,
			&booking.Status,
			&booking.TotalAmount,
			&booking.BookingDate,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return bookings, totalCount, nil
}

// GetBookingCustomerPayment retrieves a single booking customer payment by ID
func (br *BookingRepository) GetBookingCustomerPayment(ctx *gin.Context, id uint64) (*domain.BookingCustomerPayment, error) {
	var bcp domain.BookingCustomerPayment

	query := br.db.QueryBuilder.Select("*").
		From("booking_customer_payment").
		Where(sq.Eq{"booking_id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&bcp.BookingID,
		&bcp.CustomerID,
		&bcp.BookingPrice,
		&bcp.BookingStatus,
		&bcp.CheckInDate,
		&bcp.CheckOutDate,
		&bcp.BookingCreatedAt,
		&bcp.BookingUpdatedAt,
		&bcp.RoomID,
		&bcp.RoomNumber,
		&bcp.RoomTypeID,
		&bcp.RoomTypeName,
		&bcp.Floor,
		&bcp.RatePriceID,
		&bcp.CustomerFirstName,
		&bcp.CustomerSurname,
		&bcp.PaymentID,
		&bcp.PaymentStatus,
		&bcp.PaymentUpdateDate,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &bcp, nil
}

// ListBookingCustomerPayments retrieves a list of booking customer payments with pagination
func (br *BookingRepository) ListBookingCustomerPayments(ctx *gin.Context, skip, limit uint64) ([]domain.BookingCustomerPayment, uint64, error) {
	var bcps []domain.BookingCustomerPayment
	var totalCount uint64

	countQuery := br.db.QueryBuilder.Select("COUNT(*)").
		From("booking_customer_payment")
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = br.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := br.db.QueryBuilder.Select("*").
		From("booking_customer_payment").
		OrderBy("booking_id").
		Limit(limit)

	if skip > 0 {
		query = query.Offset(skip)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := br.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var bcp domain.BookingCustomerPayment
		err := rows.Scan(
			&bcp.BookingID,
			&bcp.CustomerID,
			&bcp.BookingPrice,
			&bcp.BookingStatus,
			&bcp.CheckInDate,
			&bcp.CheckOutDate,
			&bcp.BookingCreatedAt,
			&bcp.BookingUpdatedAt,
			&bcp.RoomID,
			&bcp.RoomNumber,
			&bcp.RoomTypeID,
			&bcp.RoomTypeName,
			&bcp.Floor,
			&bcp.RatePriceID,
			&bcp.CustomerFirstName,
			&bcp.CustomerSurname,
			&bcp.PaymentID,
			&bcp.PaymentStatus,
			&bcp.PaymentUpdateDate,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
		}

		bcps = append(bcps, bcp)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return bcps, totalCount, nil
}

func (br *BookingRepository) UpdateBooking(ctx *gin.Context, booking *domain.Booking) (*domain.Booking, error) {
	query := br.db.QueryBuilder.Update("bookings").
		Set("customer_id", sq.Expr("COALESCE(?, customer_id)", booking.CustomerID)).
		Set("rate_prices_id", sq.Expr("COALESCE(?, rate_prices_id)", booking.RatePriceId)).
		Set("room_id", sq.Expr("COALESCE(?, room_id)", booking.RoomID)).
		Set("room_type_id", sq.Expr("COALESCE(?, room_type_id)", booking.RoomTypeID)).
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
	slog.Debug("SQL QUERY", "query", query)

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.RatePriceId,
		&booking.RoomID,
		&booking.RoomTypeID,
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

func (br *BookingRepository) DeleteBooking(ctx *gin.Context, id uint64) error {
	query := br.db.QueryBuilder.Delete("bookings").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = br.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}