package repository

import (
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CustomerTypeRepository struct {
	db *postgres.DB
}

func NewCustomerTypeRepository(db *postgres.DB) *CustomerTypeRepository {
	return &CustomerTypeRepository{
		db,
	}
}

func (ctr *CustomerTypeRepository) CreateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error) {
	query := ctr.db.QueryBuilder.Insert("customer_types").
		Columns("name", "description").
		Values(customerType.Name, customerType.Description).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = ctr.db.QueryRow(ctx, sql, args...).Scan(
		&customerType.ID,
		&customerType.Name,
		&customerType.Description,
		&customerType.CreatedAt,
		&customerType.UpdatedAt,
	)

	if err != nil {
		if errCode := ctr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return customerType, nil
}

func (ctr *CustomerTypeRepository) GetCustomerTypeByID(ctx *gin.Context, id uint64) (*domain.CustomerType, error) {
	var customerType domain.CustomerType

	query := ctr.db.QueryBuilder.Select("*").
		From("customer_types").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = ctr.db.QueryRow(ctx, sql, args...).Scan(
		&customerType.ID,
		&customerType.Name,
		&customerType.Description,
		&customerType.CreatedAt,
		&customerType.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &customerType, nil
}

func (ctr *CustomerTypeRepository) ListCustomerTypes(ctx *gin.Context, skip, limit uint64) ([]domain.CustomerType, uint64, error) {
	var customerTypes []domain.CustomerType
	var totalCount uint64

	countQuery := ctr.db.QueryBuilder.Select("COUNT(*)").From("customer_types")
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = ctr.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := ctr.db.QueryBuilder.Select("*").
		From("customer_types").
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

	rows, err := ctr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var customerType domain.CustomerType
		err := rows.Scan(
			&customerType.ID,
			&customerType.Name,
			&customerType.Description,
			&customerType.CreatedAt,
			&customerType.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		customerTypes = append(customerTypes, customerType)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return customerTypes, totalCount, nil
}

func (ctr *CustomerTypeRepository) UpdateCustomerType(ctx *gin.Context, customerType *domain.CustomerType) (*domain.CustomerType, error) {
	query := ctr.db.QueryBuilder.Update("customer_types").
		Set("name", sq.Expr("COALESCE(?, name)", customerType.Name)).
		Set("description", sq.Expr("COALESCE(?, description)", customerType.Description)).
		Where(sq.Eq{"id": customerType.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = ctr.db.QueryRow(ctx, sql, args...).Scan(
		&customerType.ID,
		&customerType.Name,
		&customerType.Description,
		&customerType.CreatedAt,
		&customerType.UpdatedAt,
	)

	if err != nil {
		if errCode := ctr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return customerType, nil
}

func (ctr *CustomerTypeRepository) DeleteCustomerType(ctx *gin.Context, id uint64) error {
	query := ctr.db.QueryBuilder.Delete("customer_types").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = ctr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
