package main

import (
	"log"

	"github.com/hampgoodwin/GoLuca/internal/environment"
	"go.uber.org/zap"
)

func main() {
	env, err := environment.NewEnvironment(nil)
	if err != nil {
		log.Panic("failed to create new environment")
	}

	if err := env.Server.ListenAndServe(); err != nil {
		env.Log.Panic("http server failed", zap.Error(err))
	}
}
