package main

import (
	"flag"
	"net/http"

	"github.com/abelgoodwin1988/GoLuca/api"
	"github.com/abelgoodwin1988/GoLuca/internal/configloader"
	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/internal/lucalog"
	"go.uber.org/zap"
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	lucalog.Logger = logger

	if err := configloader.Load(); err != nil {
		lucalog.Logger.Fatal("failed to load config", zap.Error(err))
	}
	if err := data.CreateDB(); err != nil {
		lucalog.Logger.Fatal("failed to create new DB", zap.Error(err))
	}
	if err := data.Migrate(); err != nil {
		lucalog.Logger.Fatal("failed to migrate", zap.Error(err))
	}

	r := api.Register()

	if err := http.ListenAndServe(":3333", r); err != nil {
		lucalog.Logger.Fatal("api failure", zap.Error(err))
	}
}
