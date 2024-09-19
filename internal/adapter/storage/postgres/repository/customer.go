package repository

import (
	"context"
	"log/slog"

	"github.com/Coke3a/HotelManagement/internal/adapter/storage/postgres"
	"github.com/Coke3a/HotelManagement/internal/core/domain"
	sq "github.com/Masterminds/squirrel"
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

func (cr *CustomerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	query := cr.db.QueryBuilder.Insert("customers").
		Columns("name", "email", "phone", "address", "date_of_birth", "gender", "membership_status", "join_date", "preferences", "last_visit_date").
		Values(customer.Name, customer.Email, customer.Phone, customer.Address, customer.DateOfBirth, customer.Gender, customer.MembershipStatus, customer.JoinDate, customer.Preferences, customer.LastVisitDate).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.DateOfBirth,
		&customer.Gender,
		&customer.MembershipStatus,
		&customer.JoinDate,
		&customer.Preferences,
		&customer.LastVisitDate,
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

func (cr *CustomerRepository) GetCustomerByID(ctx context.Context, id uint64) (*domain.Customer, error) {
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
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.DateOfBirth,
		&customer.Gender,
		&customer.MembershipStatus,
		&customer.JoinDate,
		&customer.Preferences,
		&customer.LastVisitDate,
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

func (cr *CustomerRepository) ListCustomers(ctx context.Context, skip, limit uint64) ([]domain.Customer, error) {
	var customers []domain.Customer

	query := cr.db.QueryBuilder.Select("*").
		From("customers").
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

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customer domain.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.DateOfBirth,
			&customer.Gender,
			&customer.MembershipStatus,
			&customer.JoinDate,
			&customer.Preferences,
			&customer.LastVisitDate,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (cr *CustomerRepository) UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	query := cr.db.QueryBuilder.Update("customers").
		Set("name", sq.Expr("COALESCE(?, name)", customer.Name)).
		Set("email", sq.Expr("COALESCE(?, email)", customer.Email)).
		Set("phone", sq.Expr("COALESCE(?, phone)", customer.Phone)).
		Set("address", sq.Expr("COALESCE(?, address)", customer.Address)).
		Set("date_of_birth", sq.Expr("COALESCE(?, date_of_birth)", customer.DateOfBirth)).
		Set("gender", sq.Expr("COALESCE(?, gender)", customer.Gender)).
		Set("membership_status", sq.Expr("COALESCE(?, membership_status)", customer.MembershipStatus)).
		Set("preferences", sq.Expr("COALESCE(?, preferences)", customer.Preferences)).
		Set("last_visit_date", sq.Expr("COALESCE(?, last_visit_date)", customer.LastVisitDate)).
		Where(sq.Eq{"id": customer.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Debug("SQL QUERY", "query", query)

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.DateOfBirth,
		&customer.Gender,
		&customer.MembershipStatus,
		&customer.JoinDate,
		&customer.Preferences,
		&customer.LastVisitDate,
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

func (cr *CustomerRepository) DeleteCustomer(ctx context.Context, id uint64) error {
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
