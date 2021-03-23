package lucalog

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Logger is the gobally accessible logger
var Logger *zap.Logger

func Set() error {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "failed to create new zap logger")
	}

	Logger.Info("logger set and ready")
	return nil
}
