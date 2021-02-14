package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abelgoodwin1988/GoLuca/internal/configloader"
	"github.com/abelgoodwin1988/GoLuca/internal/data"
)

func main() {
	ctx := context.Background()
	if err := configloader.Load(); err != nil {
		log.Fatalf("failed to load config\n%s", err.Error())
	}
	if err := data.CreateDB(ctx); err != nil {
		log.Fatalf("failed to create new DB\n%s", err.Error())
	}
	if err := data.Migrate(ctx); err != nil {
		log.Fatalf("failed to migrate\n%s", err.Error())
	}
	run()
}

func run() {
	fmt.Println("run")
}
