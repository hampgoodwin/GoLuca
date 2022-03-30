package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/hampgoodwin/GoLuca/internal/controller"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/environment"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/router"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	env, err := environment.New(environment.Environment{}, "/etc/goluca/.env.toml")
	if err != nil {
		log.Panic("failed to create new environment")
	}

	db, err := database.NewDatabasePool(ctx, env.Config.Database.ConnectionString())
	if err != nil {
		env.Log.Error("creating new database pool", zap.Error(err))
		log.Fatal("error creating database pool on application start")
	}
	if err := database.Migrate(db); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			env.Log.Fatal("migrating", zap.Error(err))
			log.Fatal("error migrating database on application start")
		}
		env.Log.Info("no migration changes")
	}

	repository := repository.NewRepository(db)
	service := service.NewService(env.Log, repository)
	controller := controller.NewController(env.Log, service)
	s := &http.Server{
		Addr:     env.Config.HTTPAPI.AddressString(),
		ErrorLog: zap.NewStdLog(env.Log),
		Handler: router.Register(
			env.Log,
			controller.RegisterAccountRoutes,
			controller.RegisterTransactionRoutes,
		),
	}

	if err := s.ListenAndServe(); err != nil {
		env.Log.Panic("http server failed", zap.Error(err))
	}

	if err := s.Shutdown(ctx); err != nil {
		env.Log.Fatal("shutting server down", zap.Error(err))
	}
}
