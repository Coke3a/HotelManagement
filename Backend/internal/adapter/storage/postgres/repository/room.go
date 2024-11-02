package repository

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"fmt"
)

type RoomRepository struct {
	db *postgres.DB
}

func NewRoomRepository(db *postgres.DB) *RoomRepository {
	return &RoomRepository{
		db,
	}
}

func (rr *RoomRepository) CreateRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error) {
	query := rr.db.QueryBuilder.Insert("rooms").
		Columns("room_number", "type_id", "description", "status", "floor").
		Values(room.RoomNumber, room.TypeID, room.Description, int(room.Status), room.Floor).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = rr.db.QueryRow(ctx, sql, args...).Scan(
		&room.ID,
		&room.RoomNumber,
		&room.TypeID,
		&room.Description,
		&room.Status,
		&room.Floor,
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

func (rr *RoomRepository) GetRoomByID(ctx *gin.Context, id uint64) (*domain.Room, error) {
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
		&room.TypeID,
		&room.Description,
		&room.Status,
		&room.Floor,
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

func (rr *RoomRepository) ListRooms(ctx *gin.Context, skip, limit uint64) ([]domain.Room, error) {
	var rooms []domain.Room

	query := rr.db.QueryBuilder.Select("*").
		From("rooms").
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
			&room.TypeID,
			&room.Description,
			&room.Status,
			&room.Floor,
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

func (rr *RoomRepository) UpdateRoom(ctx *gin.Context, room *domain.Room) (*domain.Room, error) {
	query := rr.db.QueryBuilder.Update("rooms").
		Set("room_number", sq.Expr("COALESCE(?, room_number)", room.RoomNumber)).
		Set("type_id", sq.Expr("COALESCE(?, type_id)", room.TypeID)).
		Set("description", sq.Expr("COALESCE(?, description)", room.Description)).
		Set("status", sq.Expr("COALESCE(?, status)", room.Status)).
		Set("floor", sq.Expr("COALESCE(?, floor)", room.Floor)).
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
		&room.TypeID,
		&room.Description,
		&room.Status,
		&room.Floor,
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

func (rr *RoomRepository) DeleteRoom(ctx *gin.Context, id uint64) error {
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

func (rr *RoomRepository) GetAvailableRooms(ctx *gin.Context, checkInDate, checkOutDate time.Time) ([]domain.RoomWithRoomType, error) {
	var rooms []domain.RoomWithRoomType

	query := `
	SELECT DISTINCT
		r.id, r.room_number, r.type_id, rt.name AS room_type_name, r.description, r.status, r.floor, r.created_at, r.updated_at
	FROM
		rooms r
	JOIN room_types rt ON r.type_id = rt.id
	WHERE
		r.status = 1
		AND NOT EXISTS (
			SELECT 1
			FROM bookings b
			WHERE b.room_id = r.id
			AND b.status NOT IN (5, 6) -- Exclude canceled or completed bookings
			AND (
				(b.check_in_date < $2 AND b.check_out_date > $1)  -- Overlaps with the new booking period
			)
		)`

	rows, err := rr.db.Query(ctx, query, checkInDate, checkOutDate)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var room domain.RoomWithRoomType
		err := rows.Scan(
			&room.ID,
			&room.RoomNumber,
			&room.TypeID,
			&room.TypeName,
			&room.Description,
			&room.Status,
			&room.Floor,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return rooms, nil
}

func (rr *RoomRepository) ListRoomsWithRoomType(ctx *gin.Context, skip, limit uint64) ([]domain.RoomWithRoomType, error) {
	var rooms []domain.RoomWithRoomType

	query := `
		SELECT r.id, r.room_number, r.type_id, rt.name AS room_type_name, r.description, r.status, r.floor, r.created_at, r.updated_at
		FROM rooms r
		JOIN room_types rt ON r.type_id = rt.id
		ORDER BY r.id
		LIMIT $1 OFFSET $2
	`

	rows, err := rr.db.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var room domain.RoomWithRoomType
		err := rows.Scan(
			&room.ID,
			&room.RoomNumber,
			&room.TypeID,
			&room.TypeName,
			&room.Description,
			&room.Status,
			&room.Floor,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return rooms, nil
}
