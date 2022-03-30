package environment

import (
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/configloader"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Environment struct {
	Config config.Config
	Log    *zap.Logger
}

func New(e Environment, fp string) (Environment, error) {
	env := Environment{}
	if e != (Environment{}) {
		env = e
	}

	logger, _ := zap.NewProduction()

	defer func() { _ = logger.Sync() }()

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

	return env, nil
}
