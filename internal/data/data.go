package data

import (
	"context"
	"fmt"

	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/abelgoodwin1988/GoLuca/internal/lucalog"
	"github.com/abelgoodwin1988/GoLuca/internal/setup"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// DBPool is the app-wide accessible pgx conn pool
var DBPool *pgxpool.Pool

// CreateDB creates and puts in memory a DB
func CreateDB() {
	<-setup.C.ReadyForDBCreation()

	ctx := context.Background()
	var err error
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.Env.DBUser,
		config.Env.DBPass,
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBDB,
	)

	DBPool, err = pgxpool.Connect(ctx, connString)
	if err != nil {
		setup.C.Err <- errors.Wrap(err, "failed to create pgx connection pool")
	}

	setup.C.Mu.Lock()
	setup.C.DB.Ready = true
	setup.C.Mu.Unlock()

	lucalog.Logger.Info("db connection established")
}

// Migrate handles the db migration logic. Eventually this should be replaced with a well-tested migration tool
func Migrate() {
	<-setup.C.ReadyForMigration()
	ctx := context.Background()
	tx, err := DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		setup.C.Err <- err
	}
	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS account(
	id BIGSERIAL PRIMARY KEY,
	parent_id INT,
	name VARCHAR(255) UNIQUE,
	type SMALLINT,
	basis VARCHAR(6),
	created_at TIMESTAMP DEFAULT NOW(),
	modified_at TIMESTAMP DEFAULT NOW()
)
;`)
	if err != nil {
		setup.C.Err <- errors.Wrap(err, "failed to create account table")
	}
	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS transaction(
	id BIGSERIAL PRIMARY KEY,
	description TEXT,
	created_at TIMESTAMP DEFAULT NOW()
)
;`)
	if err != nil {
		setup.C.Err <- errors.Wrap(err, "failed to create transaction table")
	}
	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS entry(
	id BIGSERIAL PRIMARY KEY,
	transaction_id int,
	account_id int,
	amount DOUBLE PRECISION,
	created_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_transaction FOREIGN KEY(transaction_id) REFERENCES transaction(id),
	CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES account(id)
)
;`)
	if err != nil {
		setup.C.Err <- errors.Wrap(err, "failed to create entry table")
	}
	if err := tx.Commit(ctx); err != nil {
		setup.C.Err <- errors.Wrap(err, "failed to commit migration")
	}
	setup.C.Mu.Lock()
	setup.C.Migration.Ready = true
	setup.C.Mu.Unlock()
	lucalog.Logger.Info("migration successful")
}
