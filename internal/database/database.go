package database

import (
	"context"
	"fmt"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// NewDatabasePool creates a new DB
func NewDatabasePool(connString string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	var err error
	DBPool, err := pgxpool.Connect(ctx, connString)
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

// Migrate handles the db migration logic. Eventually this should be replaced with a well-tested migration tool
func Migrate(DBPool *pgxpool.Pool, log *zap.Logger) error {
	ctx := context.Background()
	tx, err := DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS account(
	id VARCHAR(36) PRIMARY KEY,
	parent_id VARCHAR(36),
	name VARCHAR(255) UNIQUE,
	type VARCHAR(64),
	basis VARCHAR(6),
	created_at TIMESTAMP DEFAULT NOW()
)
;`)
	if err != nil {
		return errors.Wrap(err, "failed to create account table")
	}
	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS transaction(
	id VARCHAR(36) PRIMARY KEY,
	description TEXT,
	created_at TIMESTAMP DEFAULT NOW()
)
;`)
	if err != nil {
		return errors.Wrap(err, "failed to create transaction table")
	}
	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS entry(
	id VARCHAR(36) PRIMARY KEY,
	transaction_id VARCHAR(36),
	debit_account VARCHAR(36),
	credit_account VARCHAR(36),
	amount_value BIGINT,
	amount_currency CHAR(3),
	created_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_transaction FOREIGN KEY(transaction_id) REFERENCES transaction(id),
	CONSTRAINT fk_debit_account FOREIGN KEY(debit_account) REFERENCES account(id),
	CONSTRAINT fk_credit_account FOREIGN KEY(credit_account) REFERENCES account(id)
)
;`)
	if err != nil {
		return errors.Wrap(err, "failed to create entry table")
	}

	_, err = tx.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS transaction_created_at_idx ON transaction (id, created_at);`)
	if err != nil {
		return errors.Wrap(err, "failed to create index on transaction table")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit migration")
	}
	log.Info("migration successful")
	return nil
}
