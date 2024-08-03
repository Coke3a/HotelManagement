package repository

import (
	"context"
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type RoomRepository struct {
	db *postgres.DB
}

func NewRoomRepository(db *postgres.DB) *RoomRepository {
	return &RoomRepository{
		db,
	}
}

func (rr *RoomRepository) CreateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	query := rr.db.QueryBuilder.Insert("rooms").
		Columns("room_number", "type", "description", "status", "floor", "capacity", "price_per_night").
		Values(room.RoomNumber, room.Type, room.Description, room.Status, room.Floor, room.Capacity, room.PricePerNight).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rr.db.QueryRow(ctx, sql, args...).Scan(
		&room.ID,
		&room.RoomNumber,
		&room.Type,
		&room.Description,
		&room.Status,
		&room.Floor,
		&room.Capacity,
		&room.PricePerNight,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if errCode := rr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return room, nil
}

func (rr *RoomRepository) GetRoomByID(ctx context.Context, id uint64) (*domain.Room, error) {
	var room domain.Room

	query := rr.db.QueryBuilder.Select("*").
		From("rooms").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rr.db.QueryRow(ctx, sql, args...).Scan(
		&room.ID,
		&room.RoomNumber,
		&room.Type,
		&room.Description,
		&room.Status,
		&room.Floor,
		&room.Capacity,
		&room.PricePerNight,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &room, nil
}

func (rr *RoomRepository) ListRooms(ctx context.Context, skip, limit uint64) ([]domain.Room, error) {
	var rooms []domain.Room

	query := rr.db.QueryBuilder.Select("*").
		From("rooms").
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
		var room domain.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomNumber,
			&room.Type,
			&room.Description,
			&room.Status,
			&room.Floor,
			&room.Capacity,
			&room.PricePerNight,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (rr *RoomRepository) UpdateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	query := rr.db.QueryBuilder.Update("rooms").
		Set("room_number", sq.Expr("COALESCE(?, room_number)", room.RoomNumber)).
		Set("type", sq.Expr("COALESCE(?, type)", room.Type)).
		Set("description", sq.Expr("COALESCE(?, description)", room.Description)).
		Set("status", sq.Expr("COALESCE(?, status)", room.Status)).
		Set("floor", sq.Expr("COALESCE(?, floor)", room.Floor)).
		Set("capacity", sq.Expr("COALESCE(?, capacity)", room.Capacity)).
		Set("price_per_night", sq.Expr("COALESCE(?, price_per_night)", room.PricePerNight)).
		Where(sq.Eq{"id": room.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rr.db.QueryRow(ctx, sql, args...).Scan(
		&room.ID,
		&room.RoomNumber,
		&room.Type,
		&room.Description,
		&room.Status,
		&room.Floor,
		&room.Capacity,
		&room.PricePerNight,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if errCode := rr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return room, nil
}

func (rr *RoomRepository) DeleteRoom(ctx context.Context, id uint64) error {
	query := rr.db.QueryBuilder.Delete("rooms").
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
