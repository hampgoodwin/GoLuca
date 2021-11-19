package environment

import (
	"net/http"

	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/configloader"
	"github.com/hampgoodwin/GoLuca/internal/controller"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Environment struct {
	Config config.Config
	Log    *zap.Logger

	// API: Server and controllers
	Server     *http.Server
	controller *controller.Controller

	// service
	service *service.Service

	// DATA: database and repo
	database   *pgxpool.Pool
	repository *repository.Repository
}

func New(e Environment, fp string) (Environment, error) {
	env := Environment{}
	if e != (Environment{}) {
		env = e
	}

	logger, _ := zap.NewProduction()

	var err error
	env.Config, err = configloader.Load(env.Config, fp)
	if err != nil {
		logger.Fatal("loading configuration", zap.Error(err))
	}

	if env.Log == nil {
		switch env.Config.Environment.Type {
		case "PROD":
			env.Log = logger
		case "STAGING":
			env.Log = logger.WithOptions(zap.AddCaller())
		case "DEV":
			env.Log, _ = zap.NewDevelopment()
			env.Log = env.Log.WithOptions(zap.AddCaller())
		case "LOCAL":
			env.Log, _ = zap.NewDevelopment()
			env.Log = env.Log.WithOptions(zap.AddCaller())
		default:
			env.Log = logger
		}
		env.Log = env.Log.With(
			zap.String("application", "goluca"),
			zap.String("environment", env.Config.Environment.Type))
	}

	// Database
	if env.database == nil {
		db, err := database.NewDatabasePool(env.Config.Database.ConnectionString())
		if err != nil {
			env.Log.Fatal("creating new database pool", zap.Error(err))
		}
		env, err = SetDatabase(env, db)
		if err != nil {
			env.Log.Fatal("setting new database", zap.Error(err))
		}

	}

	if env.repository == nil {
		env, err = SetRepository(env, env.database)
		if err != nil {
			return env, errors.Wrap(err, "setting environment repository")
		}
	}

	// Service
	if env.service == nil {
		env, err = SetService(env, env.repository)
		if err != nil {
			return env, errors.Wrap(err, "setting environment service")
		}
	}

	// Controllers
	if env.controller == nil {
		env, err = SetController(env, env.service)
		if err != nil {
			return env, errors.Wrap(err, "setting environment controller")
		}
	}

	if env.Server == nil {
		env.Server = &http.Server{
			Addr:     env.Config.HTTPAPI.AddressString(),
			ErrorLog: zap.NewStdLog(env.Log),
		}
	}
	// register routes
	env.Server.Handler = controller.Register(
		env.Log,
		env.controller.RegisterAccountRoutes,
		env.controller.RegisterTransactionRoutes,
	)

	return env, nil
}

func SetDatabase(env Environment, db *pgxpool.Pool) (Environment, error) {
	env.database = db
	return env, nil
}

func MigrateDatabase(env Environment) error {
	if err := database.Migrate(env.database, env.Log); err != nil {
		env.Log.Fatal("migrating", zap.Error(err))
	}
	return nil
}

func CloseDatabase(env Environment) {
	env.database.Close()
}

func SetRepository(env Environment, db *pgxpool.Pool) (Environment, error) {
	if db == nil {
		return env, errors.New("cannot set environment repository without database")
	}
	env.repository = repository.NewRepository(db)
	return env, nil
}

func SetService(env Environment, r *repository.Repository) (Environment, error) {
	if r == nil {
		return env, errors.New("cannot set environment service without repository")
	}
	env.service = service.NewService(env.Log, r)
	return env, nil
}

func SetController(env Environment, s *service.Service) (Environment, error) {
	if s == nil {
		return env, errors.New("cannot set environment controller without service")
	}
	env.controller = controller.NewController(env.Log, s)
	return env, nil
}

func StartHTTPServer(env Environment) func() {
	if err := env.Server.ListenAndServe(); err != nil {
		env.Log.Panic("http server failed", zap.Error(err))
	}
	shutdown := func() { _ = env.Server.Shutdown }
	return shutdown
}
