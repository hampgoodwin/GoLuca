package data

import (
	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// DB is the app-wide accessable DB
var DB *sqlx.DB

// Create creates and puts in memory a DB
func Create() error {
	var err error
	DB, err = sqlx.Open(config.Env.DBDriverName, config.Env.DBConnectionString)
	if err != nil {
		return err
	}
	return nil
}
