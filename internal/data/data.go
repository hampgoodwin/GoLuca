package data

import (
	"context"
	"fmt"

	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/jmoiron/sqlx"

	// postgres driver
	_ "github.com/lib/pq"
)

// DB is the app-wide accessable DB
var DB *sqlx.DB

// CreateDB creates and puts in memory a DB
func CreateDB(ctx context.Context) error {
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
func Migrate(ctx context.Context) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
USE goluca;
CREATE TABLE IF NOT EXISTS entry(
	id SERIAL PRIMARY KEY,
	account_id VARCHAR(255),
	amount DOUBLE PRECISION
)
;`)
	if err != nil {
		return err
	}
	fmt.Println("migration successful")
	return nil
}
