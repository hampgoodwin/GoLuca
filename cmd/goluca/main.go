package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	accountconnectv1 "github.com/hampgoodwin/GoLuca/internal/account/connectrpc/v1"
	accountrepository "github.com/hampgoodwin/GoLuca/internal/account/repository"
	accountservice "github.com/hampgoodwin/GoLuca/internal/account/service"
	"github.com/hampgoodwin/GoLuca/internal/database"
	postgresdb "github.com/hampgoodwin/GoLuca/internal/database/postgres"
	"github.com/hampgoodwin/GoLuca/internal/environment"
	inats "github.com/hampgoodwin/GoLuca/internal/event/nats"
	itrace "github.com/hampgoodwin/GoLuca/internal/trace"
	transactionconnectv1 "github.com/hampgoodwin/GoLuca/internal/transaction/connectrpc/v1"
	transactionrepository "github.com/hampgoodwin/GoLuca/internal/transaction/repository"
	transactionservice "github.com/hampgoodwin/GoLuca/internal/transaction/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Create the OTLP Tracer Provider
	tpShutdownFn, err := itrace.SetOTLPGRPCTracerProvider(ctx)
	if err != nil {
		log.Panic("failed to create otlp grpc exporter")
	}

	// Load the environment
	// The environment includes the minimum necessary dependencies to start the application
	env, err := environment.New(environment.Environment{}, "/etc/goluca/.env.toml")
	if err != nil {
		log.Panic("failed to create new environment")
	}

	// Create the postgres database pool and migrate
	db, err := postgresdb.NewDatabasePool(ctx, env.Config.Database.ConnectionString())
	if err != nil {
		env.Log.Error("creating new database pool", zap.Error(err))
		log.Fatal("error creating database pool on application start")
	}
	if err := postgresdb.Migrate(db, database.MigrationsFS); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			env.Log.Fatal("migrating", zap.Error(err))
			log.Fatal("error migrating database on application start")
		}
		env.Log.Info("no migration changes")
	}

	// Create NATS event bus, using proto encoded connection and JetStream
	env.Log.Info("starting nats", zap.Any("service_info", env.Config.NATS.URL()))
	nc, err := inats.NewNATSConn(env.Config.NATS.URL())
	if err != nil {
		env.Log.Error("nats error, shutting down", zap.Error(err))
		close(ctx, env.Log, db, nc, nil, nil, tpShutdownFn)
		log.Fatal("failed to create nats connection")
	}
	env.Log.Info("creating jetstream")
	var ncWiretap *nats.Conn
	if env.Config.NATS.Wiretap.Enable {
		env.Log.Info("starting wiretap", zap.Any("service_info", env.Config.NATS.Wiretap.URL()))
		ncWiretap, err = inats.WireTap(env.Config.NATS.Wiretap.URL())
		if err != nil {
			env.Log.Error("nats wiretap error, shutting down", zap.Error(err))
			close(ctx, env.Log, db, nc, ncWiretap, nil, tpShutdownFn)
			log.Fatal("failed to create wiretap")
		}
	}

	// create layers
	//// account
	accountRepository := accountrepository.NewRepository(db)
	accountService := accountservice.NewService(accountRepository, nc)
	accountHandler := accountconnectv1.NewHandler(accountService)
	//// transaction
	transactionRepository := transactionrepository.NewRepository(db)
	transactionService := transactionservice.NewService(transactionRepository, nc)
	transactionHandler := transactionconnectv1.NewHandler(transactionService)

	mux := http.NewServeMux()
	accountconnectv1.Register(mux, accountHandler)
	transactionconnectv1.Register(mux, transactionHandler)

	connectServer := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	connectErr := make(chan error)
	go func() {
		env.Log.Info("starting connect server", zap.String("service_info", "localhost:8080"))
		connectErr <- connectServer.ListenAndServe()
	}()

	// Handle any errors from Servers
	for {
		select {
		case err := <-connectErr:
			env.Log.Error("connect server error, shutting down", zap.Error(err))
			close(ctx, env.Log, db, nc, ncWiretap, connectServer, tpShutdownFn)
			return
		case <-ctx.Done():
			fmt.Printf("received shutdown signal: %s\n", ctx.Err())
			cancel()
			close(ctx, env.Log, db, nc, ncWiretap, connectServer, tpShutdownFn)
			return
		}
	}
}

// close cleans up the application dependencies
func close(
	ctx context.Context,
	log *zap.Logger,
	db *pgxpool.Pool,
	nc *nats.Conn,
	ncWiretap *nats.Conn,
	connectServer *http.Server,
	tpShutdownFunc func(context.Context) error,
) {
	log.Info("closing")
	// close http server
	if connectServer != nil {
		if err := connectServer.Shutdown(ctx); err != nil {
			log.Info("closing httpserver")
			_ = connectServer.Close()
		}
	}
	// disconnect from db
	if db != nil {
		log.Info("closing db")
		db.Close()
	}
	// drain nats encoded connection
	if nc != nil {
		log.Info("draining and closing nats connection")
		if err := nc.Drain(); err != nil {
			log.Error("draining and closing nats connection", zap.Error(err))
		}
	}
	// drain nats wire tap encoded connection
	if ncWiretap != nil {
		log.Info("draining and closing wiretap connection")
		if err := ncWiretap.Drain(); err != nil {
			log.Error("draining and closing wiretap connection", zap.Error(err))
		}
	}
	// shutdown tracer provider
	log.Info("shutting down tracer provider")
	if err := tpShutdownFunc(ctx); err != nil {
		log.Error("shutting down tracer provider", zap.Error(err))
	}
}
