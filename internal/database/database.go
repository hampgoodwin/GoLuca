package database

import (
	"context"
	"fmt"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// NewDatabase creates a new DB
func NewDatabase(DBUser, DBPass, DBHost, DBPort, DBDatabase string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	var err error
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		DBUser,
		DBPass,
		DBHost,
		DBPort,
		DBDatabase,
	)

	DBPool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pgx connection pool")
	}
	return DBPool, nil
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
	debit_accont VARCHAR(36),
	credit_accont VARCHAR(36),
	amount_value BIGINT,
	amount_currency CHAR(3),
	created_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_transaction FOREIGN KEY(transaction_id) REFERENCES transaction(id),
	CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES account(id)
)
;`)
	if err != nil {
		return errors.Wrap(err, "failed to create entry table")
	}
	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit migration")
	}
	log.Info("migration successful")
	return nil
}
