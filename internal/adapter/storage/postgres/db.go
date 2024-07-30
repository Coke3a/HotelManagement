package postgres

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/Coke3a/HotelManagement/internal/adapter/config"
	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import the postgres driver for golang-migrate
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	*pgxpool.Pool
	QueryBuilder *squirrel.StatementBuilderType
	url          string
}

func Connect(ctx context.Context, config *config.DB) (*DB, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		config.Connection,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &DB{
		db,
		&psql,
		url,
	}, nil
}

// ErrorCode returns the error code of the given error
func (db *DB) ErrorCode(err error) string {
	pgErr := err.(*pgconn.PgError)
	return pgErr.Code
}

// Close closes the database connection
func (db *DB) Close() {
	db.Pool.Close()
}

// Migrate runs the database migration
func (db *DB) Migrate() error {
	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithSourceInstance("iofs", driver, db.url)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil {
		// Check if the error message contains "Dirty database"
		if strings.Contains(err.Error(), "Dirty database") {
			// Force the migration version to the last successfully applied version (e.g., 1)
			// You may need to adjust this version number based on your migrations
			forceErr := migrations.Force(1)
			if forceErr != nil {
				return forceErr
			}
			// Retry the migration after forcing the version
			err = migrations.Up()
		}
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
