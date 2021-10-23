package test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/environment"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/matryer/is"
	"go.uber.org/zap"
)

type Scope struct {
	Config     config.Config
	Server     *http.Server
	Env        environment.Environment
	DB         *pgxpool.Pool
	dbDatabase string
	IS         *is.I
	HTTPClient *http.Client
	CTX        context.Context
}

func GetScope(t *testing.T) Scope {
	t.Helper()
	s, err := NewScope(t)
	if err != nil {
		s.CleanupScope()
		t.Fatal("creating new scope")
	}
	return s
}

func NewScope(t *testing.T) (Scope, error) {
	s := Scope{}
	s.CTX = context.Background()
	s.Env = environment.TestEnvironment
	s.Env.Log = zap.NewNop()
	s.IS = is.New(t)
	s.HTTPClient = &http.Client{Timeout: time.Second * 30}
	s.Env.Server = &http.Server{
		Addr:     s.Env.Config.HTTPAPI.AddressString(),
		ErrorLog: zap.NewStdLog(s.Env.Log),
	}

	var err error

	// swap out the database for a new random one
	gofakeit.Seed(0)
	s.dbDatabase = strings.ToLower(strings.Replace(gofakeit.Name(), " ", "", -1))
	s.DB, err = database.NewDatabasePool(s.Env.Config.Database.ConnectionString())
	if err != nil {
		return s, errors.Wrap(err, "opening new database connection")
	}
	s.Env, err = environment.SetDatabase(s.Env, s.DB)
	if err != nil {
		return s, err
	}

	// Create the new database with the existing database connection pool
	if err := database.CreateDatabase(s.DB, s.dbDatabase); err != nil {
		return s, errors.Wrap(err, "creating test database")
	}
	environment.CloseDatabase(s.Env)

	dbCFG := s.Env.Config.Database
	dbCFG.Database = s.dbDatabase
	s.DB, err = database.NewDatabasePool(dbCFG.ConnectionString())
	if err != nil {
		return s, errors.Wrap(err, "opening new database connection")
	}
	s.Env, err = environment.SetDatabase(s.Env, s.DB)
	if err != nil {
		return s, errors.Wrap(err, "setting environment database to test database")
	}
	// run migration on the new database
	if err := environment.MigrateDatabase(s.Env); err != nil {
		return s, errors.Wrap(err, fmt.Sprintf("migrating test database %q", s.dbDatabase))
	}

	s.Env, err = environment.SetRepository(s.Env, s.DB)
	if err != nil {
		return s, errors.Wrap(err, "setting test repository with test database")
	}

	s.Env, err = environment.New(s.Env, "")
	if err != nil {
		return s, errors.Wrap(err, "setting up new environment")
	}

	go func() {
		if err := s.Env.Server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Fatal("http api server listening")
		}
	}()

	t.Cleanup(s.CleanupScope)
	return s, nil
}

func (s *Scope) CleanupScope() {
	environment.CloseDatabase(s.Env)

	db, _ := database.NewDatabasePool(s.Env.Config.Database.ConnectionString())
	_ = database.DropDatabase(db, s.dbDatabase)

	_ = s.Env.Server.Shutdown(s.CTX)
}
