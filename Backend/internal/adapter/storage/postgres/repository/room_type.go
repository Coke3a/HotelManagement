package repository

import (
	"context"
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type RoomTypeRepository struct {
	db *postgres.DB
}

func NewRoomTypeRepository(db *postgres.DB) *RoomTypeRepository {
	return &RoomTypeRepository{
		db,
	}
}

func (rtr *RoomTypeRepository) CreateRoomType(ctx context.Context, roomType *domain.RoomType) (*domain.RoomType, error) {
	query := rtr.db.QueryBuilder.Insert("room_types").
		Columns("name", "description", "capacity", "default_price").
		Values(roomType.Name, roomType.Description, roomType.Capacity, roomType.DefaultPrice).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rtr.db.QueryRow(ctx, sql, args...).Scan(
		&roomType.ID,
		&roomType.Name,
		&roomType.Description,
		&roomType.Capacity,
		&roomType.DefaultPrice,
		&roomType.CreatedAt,
		&roomType.UpdatedAt,
	)

	if err != nil {
		if errCode := rtr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return roomType, nil
}

func (rtr *RoomTypeRepository) GetRoomTypeByID(ctx context.Context, id uint64) (*domain.RoomType, error) {
	var roomType domain.RoomType

	query := rtr.db.QueryBuilder.Select("*").
		From("room_types").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rtr.db.QueryRow(ctx, sql, args...).Scan(
		&roomType.ID,
		&roomType.Name,
		&roomType.Description,
		&roomType.Capacity,
		&roomType.DefaultPrice,
		&roomType.CreatedAt,
		&roomType.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &roomType, nil
}

func (rtr *RoomTypeRepository) ListRoomTypes(ctx context.Context, skip, limit uint64) ([]domain.RoomType, error) {
	var roomTypes []domain.RoomType

	query := rtr.db.QueryBuilder.Select("*").
		From("room_types").
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

	rows, err := rtr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roomType domain.RoomType
		err := rows.Scan(
			&roomType.ID,
			&roomType.Name,
			&roomType.Description,
			&roomType.Capacity,
			&roomType.DefaultPrice,
			&roomType.CreatedAt,
			&roomType.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		roomTypes = append(roomTypes, roomType)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roomTypes, nil
}

func (rtr *RoomTypeRepository) UpdateRoomType(ctx context.Context, roomType *domain.RoomType) (*domain.RoomType, error) {
	query := rtr.db.QueryBuilder.Update("room_types").
		Set("name", sq.Expr("COALESCE(?, name)", roomType.Name)).
		Set("description", sq.Expr("COALESCE(?, description)", roomType.Description)).
		Set("capacity", sq.Expr("COALESCE(?, capacity)", roomType.Capacity)).
		Set("default_price", sq.Expr("COALESCE(?, default_price)", roomType.DefaultPrice)).
		Where(sq.Eq{"id": roomType.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rtr.db.QueryRow(ctx, sql, args...).Scan(
		&roomType.ID,
		&roomType.Name,
		&roomType.Description,
		&roomType.Capacity,
		&roomType.DefaultPrice,
		&roomType.CreatedAt,
		&roomType.UpdatedAt,
	)

	if err != nil {
		if errCode := rtr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return roomType, nil
}

func (rtr *RoomTypeRepository) DeleteRoomType(ctx context.Context, id uint64) error {
	query := rtr.db.QueryBuilder.Delete("room_types").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = rtr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}