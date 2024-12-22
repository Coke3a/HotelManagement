package repository

import (
	"fmt"
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CustomerRepository struct {
	db *postgres.DB
}

func NewCustomerRepository(db *postgres.DB) *CustomerRepository {
	return &CustomerRepository{
		db,
	}
}

func (cr *CustomerRepository) CreateCustomer(ctx *gin.Context, customer *domain.Customer) (*domain.Customer, error) {
	query := cr.db.QueryBuilder.Insert("customers").
		Columns("firstname", "surname", "identity_number", "email", "phone", "address", "gender", "customer_type_id", "preferences").
		Values(customer.FirstName, customer.Surname, customer.IdentityNumber, customer.Email, customer.Phone, customer.Address, customer.Gender, customer.CustomerTypeID, customer.Preferences).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.Surname,
		&customer.IdentityNumber,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.Gender,
		&customer.CustomerTypeID,
		&customer.Preferences,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if errCode := cr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return customer, nil
}

func (cr *CustomerRepository) GetCustomerByID(ctx *gin.Context, id uint64) (*domain.Customer, error) {
	var customer domain.Customer

	query := cr.db.QueryBuilder.Select("*").
		From("customers").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.Surname,
		&customer.IdentityNumber,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.Gender,
		&customer.CustomerTypeID,
		&customer.Preferences,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &customer, nil
}

func (cr *CustomerRepository) ListCustomers(ctx *gin.Context, skip, limit uint64) ([]domain.Customer, uint64, error) {
	var customers []domain.Customer
	var totalCount uint64

	countQuery := cr.db.QueryBuilder.Select("COUNT(*)").From("customers")
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = cr.db.QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := cr.db.QueryBuilder.Select("*").
		From("customers").
		OrderBy("id DESC").
		Limit(limit)

	if skip > 0 {
		query = query.Offset(skip)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}
	slog.Debug("SQL QUERY", "query", query)

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var customer domain.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.FirstName,
			&customer.Surname,
			&customer.IdentityNumber,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.Gender,
			&customer.CustomerTypeID,
			&customer.Preferences,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return customers, totalCount, nil
}

func (cr *CustomerRepository) UpdateCustomer(ctx *gin.Context, customer *domain.Customer) (*domain.Customer, error) {
	query := cr.db.QueryBuilder.Update("customers").
		Set("firstname", sq.Expr("COALESCE(?, firstname)", customer.FirstName)).
		Set("surname", sq.Expr("COALESCE(?, surname)", customer.Surname)).
		Set("identity_number", sq.Expr("COALESCE(?, identity_number)", customer.IdentityNumber)).
		Set("email", sq.Expr("COALESCE(?, email)", customer.Email)).
		Set("phone", sq.Expr("COALESCE(?, phone)", customer.Phone)).
		Set("address", sq.Expr("COALESCE(?, address)", customer.Address)).
		Set("gender", sq.Expr("COALESCE(?, gender)", customer.Gender)).
		Set("customer_type_id", sq.Expr("COALESCE(?, customer_type_id)", customer.CustomerTypeID)).
		Set("preferences", sq.Expr("COALESCE(?, preferences)", customer.Preferences)).
		Where(sq.Eq{"id": customer.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.Surname,
		&customer.IdentityNumber,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.Gender,
		&customer.CustomerTypeID,
		&customer.Preferences,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if errCode := cr.db.ErrorCode(err); errCode == "23505" {
			fmt.Println(err)
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return customer, nil
}

func (cr *CustomerRepository) DeleteCustomer(ctx *gin.Context, id uint64) error {
	query := cr.db.QueryBuilder.Delete("customers").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	slog.Debug("SQL QUERY", "query", query)

	_, err = cr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}