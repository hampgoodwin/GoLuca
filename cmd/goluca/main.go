package main

import (
	"fmt"
	"net/http"

	"github.com/abelgoodwin1988/GoLuca/api"
	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/abelgoodwin1988/GoLuca/internal/configloader"
	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/internal/lucalog"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	lucalog.Logger = logger

	if err := configloader.Load(); err != nil {
		lucalog.Logger.Fatal("failed to load config", zap.Error(err))
	}
	if err := data.CreateDB(); err != nil {
		lucalog.Logger.Fatal("failed to create new DB", zap.Error(err))
	}
	defer data.DBPool.Close()

	if err := data.Migrate(); err != nil {
		lucalog.Logger.Fatal("failed to migrate", zap.Error(err))
	}

	r := api.Register()

	server := http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%s", config.Env.APIHost, config.Env.APIPort),
	}

	if err := server.ListenAndServe(); err != nil {
		lucalog.Logger.Fatal("api failure", zap.Error(err))
	}
}
