package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	inats "github.com/hampgoodwin/GoLuca/internal/event/nats"
	"github.com/nats-io/nats.go"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	nec, err := inats.WireTap(nats.DefaultURL)
	if err != nil {
		fmt.Println("wire tapping nats")
		fmt.Println(err)
		cancel()
	}
	defer func() {
		if err := nec.Drain(); err != nil {
			fmt.Println("draining nats")
			fmt.Println(err)
			cancel()
		}
	}()

	if done := <-ctx.Done(); done != struct{}{} {
		fmt.Println("shutting down")
		cancel()
	}
}
