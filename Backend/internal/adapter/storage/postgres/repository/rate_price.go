package repository

import (
	"context"
	"log/slog"
	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type RatePriceRepository struct {
	db *postgres.DB
}

func NewRatePriceRepository(db *postgres.DB) *RatePriceRepository {
	return &RatePriceRepository{
		db,
	}
}

func (rpr *RatePriceRepository) CreateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error) {
	query := rpr.db.QueryBuilder.Insert("rate_prices").
		Columns("name", "description", "price_per_night", "room_type_id").
		Values(ratePrice.Name, ratePrice.Description, ratePrice.PricePerNight, ratePrice.RoomTypeID).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rpr.db.QueryRow(ctx, sql, args...).Scan(
		&ratePrice.ID,
		&ratePrice.Name,
		&ratePrice.Description,
		&ratePrice.PricePerNight,
		&ratePrice.RoomTypeID,
		&ratePrice.CreatedAt,
		&ratePrice.UpdatedAt,
	)

	if err != nil {
		if errCode := rpr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return ratePrice, nil
}

func (rpr *RatePriceRepository) GetRatePriceByID(ctx context.Context, id uint64) (*domain.RatePrice, error) {
	var ratePrice domain.RatePrice

	query := rpr.db.QueryBuilder.Select("*").
		From("rate_prices").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rpr.db.QueryRow(ctx, sql, args...).Scan(
		&ratePrice.ID,
		&ratePrice.Name,
		&ratePrice.Description,
		&ratePrice.PricePerNight,
		&ratePrice.RoomTypeID,
		&ratePrice.CreatedAt,
		&ratePrice.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &ratePrice, nil
}

func (rpr *RatePriceRepository) ListRatePrices(ctx context.Context, skip, limit uint64) ([]domain.RatePrice, error) {
	var ratePrices []domain.RatePrice

	query := rpr.db.QueryBuilder.Select("*").
		From("rate_prices").
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

	rows, err := rpr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ratePrice domain.RatePrice
		err := rows.Scan(
			&ratePrice.ID,
			&ratePrice.Name,
			&ratePrice.Description,
			&ratePrice.PricePerNight,
			&ratePrice.RoomTypeID,
			&ratePrice.CreatedAt,
			&ratePrice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		ratePrices = append(ratePrices, ratePrice)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ratePrices, nil
}

func (rpr *RatePriceRepository) UpdateRatePrice(ctx context.Context, ratePrice *domain.RatePrice) (*domain.RatePrice, error) {
	query := rpr.db.QueryBuilder.Update("rate_prices").
		Set("name", sq.Expr("COALESCE(?, name)", ratePrice.Name)).
		Set("description", sq.Expr("COALESCE(?, description)", ratePrice.Description)).
		Set("price_per_night", sq.Expr("COALESCE(?, price_per_night)", ratePrice.PricePerNight)).
		Set("room_type_id", sq.Expr("COALESCE(?, room_type_id)", ratePrice.RoomTypeID)).
		Where(sq.Eq{"id": ratePrice.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rpr.db.QueryRow(ctx, sql, args...).Scan(
		&ratePrice.ID,
		&ratePrice.Name,
		&ratePrice.Description,
		&ratePrice.PricePerNight,
		&ratePrice.RoomTypeID,
		&ratePrice.CreatedAt,
		&ratePrice.UpdatedAt,
	)

	if err != nil {
		if errCode := rpr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return ratePrice, nil
}

func (rpr *RatePriceRepository) DeleteRatePrice(ctx context.Context, id uint64) error {
	query := rpr.db.QueryBuilder.Delete("rate_prices").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = rpr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (rpr *RatePriceRepository) GetRatePricesByRoomTypeId(ctx context.Context, roomTypeID uint64) ([]domain.RatePrice, error) {
	var ratePrices []domain.RatePrice

	query := rpr.db.QueryBuilder.Select("*").
		From("rate_prices").
		Where(sq.Eq{"room_type_id": roomTypeID}).
		OrderBy("id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := rpr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ratePrice domain.RatePrice
		err := rows.Scan(
			&ratePrice.ID,
			&ratePrice.Name,
			&ratePrice.Description,
			&ratePrice.PricePerNight,
			&ratePrice.RoomTypeID,
			&ratePrice.CreatedAt,
			&ratePrice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		ratePrices = append(ratePrices, ratePrice)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ratePrices, nil
}

func (rpr *RatePriceRepository) GetRatePricesByRoomId(ctx context.Context, roomID uint64) ([]domain.RatePrice, error) {
    var ratePrices []domain.RatePrice

    query := `
        SELECT rp.*
        FROM rate_prices rp
        JOIN rooms r ON r.type_id = rp.room_type_id
        WHERE r.id = $1
    `

    rows, err := rpr.db.Query(ctx, query, roomID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var ratePrice domain.RatePrice
        err := rows.Scan(
            &ratePrice.ID,
            &ratePrice.Name,
            &ratePrice.Description,
            &ratePrice.PricePerNight,
            &ratePrice.RoomTypeID,
            &ratePrice.CreatedAt,
            &ratePrice.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        ratePrices = append(ratePrices, ratePrice)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return ratePrices, nil
}