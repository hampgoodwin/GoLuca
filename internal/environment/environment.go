package environment

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/api"
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/configloader"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Environment struct {
	Config *config.Config
	Log    *zap.Logger

	// API: Server and controllers
	Server     *http.Server
	controller *api.Controller

	// service
	service *service.Service

	// DATA: database and repo
	database   *pgxpool.Pool
	repository *repository.Repository
}

func NewEnvironment(e *Environment) (*Environment, error) {
	env := &Environment{}
	if e != nil {
		env = e
	}

	logger, _ := zap.NewProduction()

	var err error
	env.Config, err = configloader.Load(env.Config)
	if err != nil {
		logger.Fatal("loading configuration", zap.Error(err))
	}

	switch env.Config.EnvType {
	case "PROD":
		env.Log = logger

	case "STAGING":
		env.Log = logger.WithOptions(zap.AddCaller())

	case "DEV":
		env.Log, _ = zap.NewDevelopment()
		env.Log = env.Log.WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel))

	case "LOCAL":
		env.Log, _ = zap.NewDevelopment()
		env.Log = env.Log.WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel))
	}
	env.Log = env.Log.With(zap.String("application", "goluca"), zap.String("environment", env.Config.EnvType))

	// Database
	if env.database == nil {
		c := env.Config
		db, err := database.NewDatabase(c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBDatabase)
		if err != nil {
			env.Log.Fatal("creating new db", zap.Error(err))
		}
		env.database = db
		if err := database.Migrate(env.database, env.Log); err != nil {
			env.Log.Fatal("migrating", zap.Error(err))
		}
	}
	if env.repository == nil {
		env.repository = repository.NewRepository(env.database)
	}

	// Service
	if env.service == nil {
		env.service = service.NewService(env.Log, env.repository)
	}

	// Controllers
	if env.controller == nil {
		env.controller = api.NewController(env.Log, env.service)
	}
	// register routes
	api.Register(
		env.controller.RegisterHealthRoutes,
		env.controller.RegisterAccountRoutes,
		env.controller.RegisterTransactionRoutes,
		env.controller.RegisterEntryRoutes,
	)

	if env.Server == nil {
		env.Server = &http.Server{
			Addr:        fmt.Sprintf("%s:%s", env.Config.APIHost, env.Config.APIPort),
			ErrorLog:    zap.NewStdLog(env.Log),
			ReadTimeout: time.Second * 1,
		}
	}

	return env, nil
}
