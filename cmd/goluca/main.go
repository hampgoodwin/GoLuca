package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/environment"
	grpccontroller "github.com/hampgoodwin/GoLuca/internal/grpc/v1/controller"
	grpcrouter "github.com/hampgoodwin/GoLuca/internal/grpc/v1/router"
	httpcontroller "github.com/hampgoodwin/GoLuca/internal/http/v0/controller"
	httprouter "github.com/hampgoodwin/GoLuca/internal/http/v0/router"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	itrace "github.com/hampgoodwin/GoLuca/internal/trace"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// tpShutdownFn, err := itrace.NewStdOutExporter(f)
	// if err != nil {
	// 	log.Panic("failed to create otlp grpc exporter")
	// }
	// defer func() {
	// 	if err := tpShutdownFn(ctx); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	tpShutdownFn, err := itrace.SetOTLPGRPCTracerProvider(ctx)
	if err != nil {
		log.Panic("failed to create otlp grpc exporter")
	}
	defer func() {
		if err := tpShutdownFn(ctx); err != nil {
			log.Fatal(err)
		}
	}()

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

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_zap.UnaryServerInterceptor(env.Log),
	))
	grpcrouter.Register(grpcServer, grpcController)

	grpcErr := make(chan error)
	go func() {
		env.Log.Info("starting grpc server", zap.Any("service_info", grpcServer.GetServiceInfo()))
		grpcErr <- grpcServer.Serve(lis)
	}()

	httpErr := make(chan error)
	go func() {
		env.Log.Info("starting http server", zap.String("service_info", httpServer.Addr))
		httpErr <- httpServer.ListenAndServe()
	}()

	for {
		select {
		case err := <-grpcErr:
			env.Log.Error("grpc server error, shutting down", zap.Error(err))
			close(ctx, db, grpcServer, httpServer)
		case err := <-httpErr:
			env.Log.Error("http server error, shutting down", zap.Error(err))
			close(ctx, db, grpcServer, httpServer)
		}
	}
}

func close(ctx context.Context, db *pgxpool.Pool, grpcServer *grpc.Server, httpServer *http.Server) {
	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(ctx); err != nil {
		httpServer.Close()
	}
	db.Close()
}
