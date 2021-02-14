package data

import (
	"context"

	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/jmoiron/sqlx"

	// postgres driver
	_ "github.com/lib/pq"
)

// DB is the app-wide accessable DB
var DB *sqlx.DB

// CreateDB creates and puts in memory a DB
func CreateDB() error {
	var err error
	DB, err = sqlx.Open(config.Env.DBDriverName, config.Env.DBConnString)
	if err != nil {
		return err
	}
	// test the connection
	ctx := context.Background()
	c, err := DB.Conn(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	return nil
}
