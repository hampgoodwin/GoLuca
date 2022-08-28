package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/environment"
	grpccontroller "github.com/hampgoodwin/GoLuca/internal/grpc/v1/controller"
	grpcrouter "github.com/hampgoodwin/GoLuca/internal/grpc/v1/router"
	httpcontroller "github.com/hampgoodwin/GoLuca/internal/http/v0/controller"
	httprouter "github.com/hampgoodwin/GoLuca/internal/http/v0/router"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	httpController := httpcontroller.NewController(env.Log, service)
	httpServer := &http.Server{
		Addr:     env.Config.HTTPAPI.AddressString(),
		ErrorLog: zap.NewStdLog(env.Log),
		Handler: httprouter.Register(
			env.Log,
			httpController.RegisterAccountRoutes,
			httpController.RegisterTransactionRoutes,
		),
	}

	grpcController := grpccontroller.NewController(env.Log, service)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcrouter.Register(grpcServer, grpcController)

	if err := grpcServer.Serve(lis); err != nil {
		env.Log.Panic("grpc server failed", zap.Error(err))
	}

	if err := httpServer.ListenAndServe(); err != nil {
		env.Log.Panic("http server failed", zap.Error(err))
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		env.Log.Fatal("shutting server down", zap.Error(err))
	}
}
