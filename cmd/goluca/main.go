package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abelgoodwin1988/GoLuca/api"
	"github.com/abelgoodwin1988/GoLuca/internal/configloader"
	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/internal/lucalog"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	lucalog.Logger = logger

	if err := configloader.Load(); err != nil {
		log.Fatalf("failed to load config\n%s", err.Error())
	}
	if err := data.CreateDB(); err != nil {
		log.Fatalf("failed to create new DB\n%s", err.Error())
	}
	if err := data.Migrate(); err != nil {
		log.Fatalf("failed to migrate\n%s", err.Error())
	}

	r := api.Register()
	if err := http.ListenAndServe(":3333", r); err != nil {
		log.Fatal("api failure")
	}
}

func run() {
	fmt.Println("run")
}
