package main

import (
	"log"

	"github.com/hampgoodwin/GoLuca/internal/environment"
)

func main() {
	env, err := environment.New(environment.Environment{}, "/etc/goluca/.env.toml")
	if err != nil {
		log.Panic("failed to create new environment")
	}

	shutdown := environment.StartHTTPServer(env)
	defer shutdown()
}
