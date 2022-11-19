package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/hampgoodwin/errors"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	pgxtrace "github.com/hampgoodwin/GoLuca/internal/trace/pgx"
)

//go:embed migrations
var migrations embed.FS

// NewDatabasePool creates a new DB
func NewDatabasePool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	var err error
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	config.ConnConfig.Tracer = &pgxtrace.Tracer{}
	DBPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pgx connection pool")
	}
	return DBPool, nil
}

func CreateDatabase(conn *pgxpool.Pool, database string) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s;", database))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating database %q", database))
	}
	return nil
}

func DropDatabase(conn *pgxpool.Pool, database string) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s;", database))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("executing drop database %q", database))
	}
	return nil
}

// Migrate handles the db migration logic.
func Migrate(conn *pgxpool.Pool) error {
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return errors.Wrap(err, "opening fs for migrations")
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, conn.Config().ConnString())
	if err != nil {
		return errors.Wrap(err, "migrating database")
	}
	if err := m.Up(); err != nil {
		return errors.Wrap(err, "running migrations")
	}

	sErr, dbErr := m.Close()
	if sErr != nil {
		err = errors.Wrap(sErr, "closing migrator connection")
	}
	if dbErr != nil {
		err = errors.Wrap(err, "closing migrator connection")
	}
	if err != nil {
		return err
	}
	return nil
}
