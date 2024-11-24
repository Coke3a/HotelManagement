package repository

import (
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) CreateUser(ctx *gin.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Insert("users").
		Columns("username", "password", "role", "rank", "hire_date", "last_login", "status").
		Values(user.UserName, user.Password, user.Role, user.Rank, user.HireDate, user.LastLogin, user.Status).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.Role,
		&user.Rank,
		&user.HireDate,
		&user.LastLogin,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByID(ctx *gin.Context, id uint64) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.Role,
		&user.Rank,
		&user.HireDate,
		&user.LastLogin,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByUserName(ctx *gin.Context, userName string) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"username": userName}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.Role,
		&user.Rank,
		&user.HireDate,
		&user.LastLogin,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) ListUsers(ctx *gin.Context, skip, limit uint64) ([]domain.User, uint64, error) {
	var users []domain.User
	var totalCount uint64

	countQuery := ur.db.QueryBuilder.Select("COUNT(*)").From("users")
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}

	err = ur.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// This is a separate query that gets the paginated results
	// Only this query uses skip and limit
	query := ur.db.QueryBuilder.Select("*").
		From("users").
		OrderBy("id").
		Limit(limit)

	if skip > 0 {
			query = query.Offset(skip)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.Password,
			&user.Role,
			&user.Rank,
			&user.HireDate,
			&user.LastLogin,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return users, totalCount, nil
}

func (ur *UserRepository) UpdateUser(ctx *gin.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Update("users").
		Set("username", sq.Expr("COALESCE(?, username)", user.UserName)).
		Set("role", sq.Expr("COALESCE(?, role)", user.Role)).
		Set("rank", sq.Expr("COALESCE(?, rank)", user.Rank)).
		Set("hire_date", sq.Expr("COALESCE(?, hire_date)", user.HireDate)).
		Set("last_login", sq.Expr("COALESCE(?, last_login)", user.LastLogin)).
		Set("status", sq.Expr("COALESCE(?, status)", user.Status)).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.Role,
		&user.Rank,
		&user.HireDate,
		&user.LastLogin,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) DeleteUser(ctx *gin.Context, id uint64) error {
	query := ur.db.QueryBuilder.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
