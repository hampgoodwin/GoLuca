package postgres

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	pgxtrace "github.com/hampgoodwin/GoLuca/internal/trace/pgx"
)

// NewDatabasePool creates a new DB
func NewDatabasePool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var err error
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	config.ConnConfig.Tracer = &pgxtrace.Tracer{}
	DBPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx connection pool: %w", err)
	}
	return DBPool, nil
}

func CreateDatabase(conn *pgxpool.Pool, database string) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s;", database))
	if err != nil {
		return fmt.Errorf("creating database %q: %w", database, err)
	}
	return nil
}

func DropDatabase(conn *pgxpool.Pool, database string) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s;", database))
	if err != nil {
		return fmt.Errorf("executing drop database %q: %w", database, err)
	}
	return nil
}

// Migrate handles the db migration logic.
func Migrate(conn *pgxpool.Pool, migrationsFS embed.FS) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("opening fs for migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, conn.Config().ConnString())
	if err != nil {
		return fmt.Errorf("migrating database: %w", err)
	}
	if err := m.Up(); err != nil {
		return fmt.Errorf("running migrations, %w", err)
	}

	sErr, dbErr := m.Close()
	if sErr != nil {
		err = fmt.Errorf("closing migrator connection, %w", err)
	}
	if dbErr != nil {
		err = fmt.Errorf("closing migrator connection, %w", err)
	}
	if err != nil {
		return err
	}
	return nil
}
