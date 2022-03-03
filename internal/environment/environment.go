package environment

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/configloader"
	"github.com/hampgoodwin/GoLuca/internal/controller"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/hampgoodwin/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Environment struct {
	Config config.Config
	Log    *zap.Logger

	// API: HTTPServer and controllers
	HTTPMux    *chi.Mux
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

	defer logger.Sync()

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
			logger.Core().Enabled(zap.ErrorLevel)
			env.Log = logger.WithOptions(zap.AddCaller())
		case "LOCAL":
			logger.Core().Enabled(zap.ErrorLevel)
			env.Log = logger.WithOptions(zap.AddCaller())
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
		if err := MigrateDatabase(env); err != nil {
			env.Log.Fatal("migrating database", zap.Error(err))
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

	// register routes
	if env.HTTPMux == nil {
		env.HTTPMux = controller.Register(
			env.Log,
			env.controller.RegisterAccountRoutes,
			env.controller.RegisterTransactionRoutes,
		)
	}

	return env, nil
}

func SetDatabase(env Environment, db *pgxpool.Pool) (Environment, error) {
	env.database = db
	return env, nil
}

func MigrateDatabase(env Environment) error {
	if err := database.Migrate(env.database); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			env.Log.Fatal("migrating", zap.Error(err))
		}
		env.Log.Info("no migration changes")
		return nil
	}
	env.Log.Info("migration successful")
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

func NewHTTPServer(env Environment) *http.Server {
	s := &http.Server{
		Addr:     env.Config.HTTPAPI.AddressString(),
		ErrorLog: zap.NewStdLog(env.Log),
		Handler:  env.HTTPMux,
	}
	return s
}
