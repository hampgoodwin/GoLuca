package main

import (
	"context"
	"log"
	"net/http"

	"github.com/hampgoodwin/GoLuca/internal/environment"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	env, err := environment.New(environment.Environment{}, "/etc/goluca/.env.toml")
	if err != nil {
		log.Panic("failed to create new environment")
	}

	// s := environment.NewHTTPServer(env)
	s := &http.Server{
		Addr:     env.Config.HTTPAPI.AddressString(),
		ErrorLog: zap.NewStdLog(env.Log),
		Handler:  env.HTTPMux,
	}

	if err := s.ListenAndServe(); err != nil {
		env.Log.Panic("http server failed", zap.Error(err))
	}

	if err := s.Shutdown(ctx); err != nil {
		env.Log.Fatal("shutting server down", zap.Error(err))
	}
}
