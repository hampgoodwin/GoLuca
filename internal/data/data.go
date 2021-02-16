package data

import (
	"context"
	"fmt"

	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	// postgres driver
	_ "github.com/lib/pq"
)

// DB is the app-wide accessable DB
var DB *sqlx.DB

// CreateDB creates and puts in memory a DB
func CreateDB() error {
	ctx := context.Background()
	var err error
	DB, err = sqlx.Open(config.Env.DBDriverName, config.Env.DBConnString)
	if err != nil {
		return err
	}
	// test the connection
	c, err := DB.Conn(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	return nil
}

// Migrate handles the db migration logic. Eventually this should be replaced with a well-tested migration tool
func Migrate() error {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
CREATE TABLE IF NOT EXISTS account(
	id SERIAL PRIMARY KEY,
	parent_id INT,
	name VARCHAR(255) UNIQUE,
	type SMALLINT,
	basis VARCHAR(6)
)
;`)
	if err != nil {
		return errors.Wrap(err, "failed to create account table")
	}
	_, err = tx.Exec(`
CREATE TABLE IF NOT EXISTS transaction(
	id SERIAL PRIMARY KEY,
	description TEXT
)
;`)
	if err != nil {
		return errors.Wrap(err, "failed to create transaction table")
	}
	_, err = tx.Exec(`
CREATE TABLE IF NOT EXISTS entry(
	id SERIAL PRIMARY KEY,
	transaction_id int,
	account_id int,
	amount DOUBLE PRECISION,
	CONSTRAINT fk_transaction FOREIGN KEY(transaction_id) REFERENCES transaction(id),
	CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES account(id)
)
;`)
	if err != nil {
		return errors.Wrap(err, "failed to create entry table")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit migration")
	}
	fmt.Println("migration successful")
	return nil
}
