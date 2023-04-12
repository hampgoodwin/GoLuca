package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	modelv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/model/v1"
	"github.com/hampgoodwin/GoLuca/internal/event"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	nc, _ := nats.Connect(nats.DefaultURL)
	defer func() {
		err := nc.Drain()
		if err != nil {
			fmt.Printf("draining nats connection %v", err)
		}
	}()

	_, err := nc.Subscribe(">", messageHandler)
	if err != nil {
		panic(fmt.Errorf("subscribing to all subjects: %w", err))
	}

	if done := <-ctx.Done(); done != struct{}{} {
		fmt.Println("shutting down")
		cancel()
	}
}

func messageHandler(msg *nats.Msg) {
	switch msg.Subject {
	case event.SubjectAccountCreated:
		account := &modelv1.Account{}
		err := proto.Unmarshal(msg.Data, account)
		if err != nil {
			fmt.Printf("error unmarshaling message on subject %q", event.SubjectAccountCreated)
		}
		fmt.Printf("received %q\n%v\n", event.SubjectAccountCreated, account)
	case event.SubjectTransactionCreated:
		transaction := &modelv1.Transaction{}
		err := proto.Unmarshal(msg.Data, transaction)
		if err != nil {
			fmt.Printf("error unmarshaling message on subject %q", event.SubjectTransactionCreated)
		}
		fmt.Printf("received %q\n%v\n", event.SubjectTransactionCreated, transaction)
	default:
		fmt.Printf("unhandled event received on subject %q\n", msg.Sub.Subject)
	}
}
