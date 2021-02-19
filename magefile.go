//+build mage

package main

import (
	"context"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Up creates all docker-compose services
func Up() error {
	ctx := context.Background()
	mg.CtxDeps(ctx, DBUp)
	return nil
}

// DBUp creates a postgres instance in a docker container with which to work in for local development
func DBUp(ctx context.Context) error {
	if err := sh.Run("docker-compose", "-f", "build/package/docker-compose.yml", "up", "-d", "db"); err != nil {
		return err
	}
	return nil
}

// Down removes all docker-compose resources
func DCDown(ctx context.Context) error {
	if err := sh.Run("docker-compose", "-f", "build/package/docker-compose.yml", "down"); err != nil {
		return err
	}
	return nil
}

func Test(ctx context.Context) error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "golangci-lint", "run")
	if err != nil {
		return err
	}
	_, err = sh.Exec(nil, os.Stdout, os.Stderr, "go", "test", "./...", "-v", "--bench", ".", "--benchmem", "--covermode=count")
	if err != nil {
		return err
	}
	return nil
}
